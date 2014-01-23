package sidewalk

import (
	"os"
	"strconv"
	"time"
)

/**
 * Used to synchronize all agents and keep them in lock-step with each other,
 * the injectors, the printer and the drunk.
 */
func Clock(in chan Pedestrian, agentClock chan bool, agentDone chan bool,
	injs [](chan bool), printer chan bool, steps int, controller chan ClockRequest,
	drunk chan bool) {
	agents := 0
	wait := 500
	<-controller
	for step := 0; step <= steps; step++ {
		select {
		case cmd := <-controller:
			switch cmd.Cmd {
			case CmdPause:
				paused := true
				for paused {
					cmd := <-controller
					switch cmd.Cmd {
					case CmdResume:
						paused = false
					case CmdExit:
						// TODO
					case CmdTime:
						wait = cmd.Arg[0]
					default:
						// Do nothing
					}
				}
			case CmdExit:
				// TODO
			case CmdTime:
				wait = cmd.Arg[0]
			default:
				// Do nothing
			}
		default:
			// do nothing
		}
		for i := 0; i < agents; i++ {
			// tick the clock for agents
			agentClock <- true
		}
		agentRemove := 0
		for i := 0; i < agents; i++ {
			// get response from agents
			if !<-agentDone {
				agentRemove++
			}
		}
		agents -= agentRemove
		drunk <- true
		for _, inj := range injs {
			// tick every injector
			inj <- true
		}
		for i := 0; i < len(injs); i++ {
			// as many times as there are injectors check input
			agent := <-in
			if agent.Ok {
				agents++
			}
		}
		// Now everone has moved and the printer can print
		printer <- true
		<-printer
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}
}

/**
 * Injects players a specific cell every given step (default = 4).
 */
func Injector(clock chan bool, cell chan Msg, out chan Pedestrian, d Direction,
	ctrl chan int) {
	step := 4
	n := 0
	rsp := make(chan bool)
	for {
		select {
		case i := <-ctrl:
			step = i
		case <-clock:
			var agent Pedestrian
			if n%step == 0 { // create a new player
				agent = Pedestrian{d, true}
				cell <- Msg{rsp, agent, Regular}
				if !<-rsp {
					agent.Ok = false
				}
			}
			out <- agent
			n++
		}
	}
}

/**
 * Prints the sidewalk every timestep. The printer expects that
 * ANSI escape characters work in the terminal, as these are used to
 * overwrite the board every step.
 */
func Printer(in chan Coordinate, drunk chan Coordinate, clock chan bool, x int, y int) {
	coords := makeCoords(x, y)
	drunkCoord := Coordinate{-1, -1}
	step := 0
	for {
		select {
		case c := <-in:
			coords[c.Y][c.X] = true // An agent is here in this turn
		case c := <-drunk:
			drunkCoord = c
		case <-clock:
			os.Stdout.WriteString("\033[s\033[0;0H\r\nStep: " + strconv.Itoa(step) + "\r\n")
			hline := ""
			for i := 0; i < x+2; i++ {
				hline += "="
			}
			hline += "\r\n"
			os.Stdout.WriteString(hline)
			for i, _ := range coords {
				line := "|"
				for j, _ := range coords[len(coords)-i-1] {
					switch {
					case drunkCoord.X == j && drunkCoord.Y == (len(coords)-i-1):
						line += "@"
					case coords[len(coords)-i-1][j]:
						line += "o"
					default:
						line += " "
					}
				}
				line += "|\n"
				os.Stdout.WriteString(line)

			}
			os.Stdout.WriteString(hline)
			os.Stdout.WriteString("o = pedestrian, @ = drunk\r\n")
			os.Stdout.WriteString("\033[u")
			setCoords(false, &coords)
			step++
			clock <- true
		}
	}
}

/**
 * Represents a single cell on the sidewalk. The cells handle moving
 * pedestrians around. Whenever a cell holds a pedestrian it will try to move
 * it every timestep, until it suceeds.
 */
func Cell(in chan Msg, out [](chan Msg), printer chan Coordinate, coord Coordinate,
	agentClock chan bool, agentDone chan bool, drunk chan DrunkRequest) {
	var agent Pedestrian
	oldAgent := false
	rsp := make(chan bool)
	isDrunk := false
	for {
		done := false
		if !oldAgent {
			msg := <-in
			agent = msg.Agent
			msg.Rsp <- true
			printer <- coord
			if msg.Status == GotDrunk {
				// this cell is now permanently occupied by a drunk
				isDrunk = true
				for isDrunk {
					msg := <-in
					msg.Rsp <- false
					if msg.Status == GotSober {
						isDrunk = false
						done = true
					}
				}
			}
		}
		for !done {
			select {
			case <-agentClock:
				d := directions(agent.Dirc)
				tries := 0
				drunk <- DrunkRequest{coord, d[tries], rsp}
				if <-rsp {
					// Do not go the original direction if a drunk is nearby
					// that way. We hate drunkards!
					tries++
				}

				for tries < 3 && !done {
					if out[d[tries]] == nil {
						// Agent walked off the board
						agentDone <- false
						oldAgent = false
						done = true
					} else {
						select {
						case out[d[tries]] <- Msg{rsp, agent, Regular}:
							if <-rsp {
								oldAgent = false
								done = true
								agentDone <- true
							} else {
								tries++
							}
						case msg := <-in:
							msg.Rsp <- false
						}
					}
				}
				if !done {
					if !oldAgent {
						oldAgent = true
					} else {
						printer <- coord
					}
					agentDone <- true
				}
			case msg := <-in:
				msg.Rsp <- false
			}
		}
	}
}

/**
 * Controls the drunk (insertion and movement).
 * Cells query for the position of the drunk.
 * The process is also in charge of moving the drunk around.
 */
func Drunkard(in chan DrunkRequest, cells [][](chan Msg), ctrl chan Coordinate,
	printer chan Coordinate, clock chan bool, x int, y int) {
	rsp := make(chan bool)
	drunk := Coordinate{-1, -1}
	for {
		select {
		case req := <-in:
			req.Rsp <- drunkNearby(drunk, req.ReqCoord, x, y)
		case coord := <-ctrl:
			if !validCoord(drunk, x, y) {
				done := false
				for !done {
					select {
					case cells[coord.Y][coord.X] <- Msg{rsp, Pedestrian{Up, false}, GotDrunk}:
						done = true
					case req := <-in:
						req.Rsp <- drunkNearby(drunk, req.ReqCoord, x, y)
					}
				}
				if <-rsp {
					drunk = coord
					printer <- coord
					ctrl <- coord
				} else {
					// Indicates that something went wrong
					ctrl <- Coordinate{-1, -1}
				}
			} else {
				ctrl <- Coordinate{-1, -1}
			}
		case <-clock:
			if validCoord(drunk, x, y) {
				randomCoord := randomDirc(drunk, x, y)
				if validCoord(randomCoord, x, y) {
					done := false
					for !done {
						select {
						case cells[randomCoord.Y][randomCoord.X] <- Msg{rsp, Pedestrian{Up, false}, GotDrunk}:
							done = true
						case req := <-in:
							req.Rsp <- drunkNearby(drunk, req.ReqCoord, x, y)
						}
					}
					if <-rsp {
						done = false
						for !done {
							select {
							case cells[drunk.Y][drunk.X] <- Msg{rsp, Pedestrian{Up, false}, GotSober}:
								<-rsp
								done = true
							case req := <-in:
								req.Rsp <- drunkNearby(drunk, req.ReqCoord, x, y)
							}
						}
						drunk = randomCoord
						printer <- randomCoord
					}
				} else {
					cells[drunk.Y][drunk.X] <- Msg{rsp, Pedestrian{Up, false}, GotSober}
					<-rsp
					drunk = Coordinate{-1, -1}
					printer <- drunk
				}
			}
		}
	}
}

/**
 * Creates a rectangular sidewalk with the given dimensions, the simulation
 * runs the specified number of steps. The injectors are added at the given
 * coordinates, and with the given directions.
 */
func Sidewalk(x int, y int, injs []Coordinate, ds []Direction, steps int) {
	agentClock := make(chan bool)
	agentDone := make(chan bool)
	cells := makeMsgChs(x, y)
	clockCtrl := make(chan ClockRequest)
	clockIn := make(chan Pedestrian)
	drunkClock := make(chan bool)
	drunkCtrl := make(chan Coordinate)
	drunkPrinter := make(chan Coordinate)
	drunkReq := make(chan DrunkRequest)
	injChs := make([](chan bool), len(injs))
	injCtrl := make([](chan int), len(injs))
	printer := make(chan Coordinate)
	printerClock := make(chan bool)
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			out := getOutChs(j, i, x, y, cells)
			go Cell(cells[i][j], out, printer, Coordinate{j, i}, agentClock,
				agentDone, drunkReq)
		}
	}
	for i, _ := range injs {
		injChs[i] = make(chan bool)
		injCtrl[i] = make(chan int)
		injx := injs[i].X
		injy := injs[i].Y
		go Injector(injChs[i], cells[injy][injx], clockIn, ds[i], injCtrl[i])
	}
	go Printer(printer, drunkPrinter, printerClock, x, y)
	go Drunkard(drunkReq, cells, drunkCtrl, drunkPrinter, drunkClock, x, y)
	go Clock(clockIn, agentClock, agentDone, injChs, printerClock, steps, clockCtrl, drunkClock)
	Controller(clockCtrl, injCtrl, drunkCtrl, x, y)
}
