package sidewalk

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
 * Used to synchronize all agents and keep them in lock-step with each other,
 * the injectors and the printer
 */
func Clock(in chan Pedestrian, agentClock chan bool, agentDone chan bool,
	injs [](chan bool), printer chan bool, steps int, controller chan bool) {
	agents := 0
	<-controller
	for step := 0; step <= steps; step++ {
		select {
		case play := <-controller:
			if !play {
				for !play {
					// wait to be unpaused
					play = <-controller
				}
			}
		default:
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
		time.Sleep(500 * time.Millisecond)

		// Now everone has moved and the printer can print
		printer <- true
		<-printer
	}
}

/**
 * Injects players a specific cell every given step (default = 4)
 */
func Injector(clock chan bool, cell chan Msg, out chan Pedestrian, d Direction, ctrl chan int) {
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
				agent = Pedestrian{d, true, false}
				cell <- Msg{rsp, agent}
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
func Printer(in chan Coordinate, clock chan bool, x int, y int) {
	coords := makeCoords(x, y)
	step := 0
	for {
		select {
		case c := <-in:
			coords[c.Y][c.X] = true // An agent is here in this turn
		case <-clock:
			os.Stdout.WriteString("\033[s\033[0;0HStep " + strconv.Itoa(step) + "\r\n")
			hline := ""
			for i := 0; i < x+2; i++ {
				hline += "="
			}
			hline += "\r\n"
			os.Stdout.WriteString(hline)
			for i, _ := range coords {
				line := "|"
				for j, _ := range coords[len(coords)-i-1] {
					if coords[len(coords)-i-1][j] {
						line += "o"
					} else {
						line += " "
					}
				}
				line += "|\n"
				os.Stdout.WriteString(line)

			}
			os.Stdout.WriteString(hline)
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
	agentClock chan bool, agentDone chan bool) {
	var agent Pedestrian
	oldAgent := false
	rsp := make(chan bool)
	for {
		if !oldAgent {
			msg := <-in
			agent = msg.Agent
			msg.Rsp <- true
			printer <- coord
		}
		done := false
		for !done {
			select {
			case <-agentClock:
				d := directions(agent.Dirc)
				tries := 0

				for tries < 3 && !done {
					if out[d[tries]] == nil {
						// Agent walked off the board
						agentDone <- false
						oldAgent = false
						done = true
					} else {
						select {
						case out[d[tries]] <- Msg{rsp, agent}:
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
 * Creates a rectangular sidewalk with the given dimensions, the simulation
 * runs the specified number of steps. The injectors are added at the given
 * coordinates, and with the given directions.
 */
func Sidewalk(x int, y int, injs []Coordinate, ds []Direction, steps int) {
	clockCtrl := make(chan bool)
	injChs := make([](chan bool), len(injs))
	injCtrl := make([](chan int), len(injs))
	clockIn := make(chan Pedestrian)
	printerClock := make(chan bool)
	agentClock := make(chan bool)
	agentDone := make(chan bool)
	printer := make(chan Coordinate)
	cells := makeMsgChs(x, y)
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			out := getOutChs(j, i, x, y, cells)
			go Cell(cells[i][j], out, printer, Coordinate{j, i}, agentClock, agentDone)
		}
	}
	for i, _ := range injs {
		injChs[i] = make(chan bool)
		injCtrl[i] = make(chan int)
		injx := injs[i].X
		injy := injs[i].Y
		go Injector(injChs[i], cells[injy][injx], clockIn, ds[i], injCtrl[i])
	}
	go Printer(printer, printerClock, x, y)
	go Clock(clockIn, agentClock, agentDone, injChs, printerClock, steps, clockCtrl)
	Controller(clockCtrl, injCtrl, y)
}

/**
 * The shell used for interaction with the player. The shell expects ANSI
 * escape characters to work in the terminal, as these are used to overwrite
 * old output.
 */
func Controller(clock chan bool, injs [](chan int), y int) {
	os.Stdout.WriteString("\033[0;0H\033[J\033[" + strconv.Itoa(y+4) + ";0H")
	clock <- true
	for {
		os.Stdout.WriteString("\033[Kcmd> ")
		stdin := bufio.NewReader(os.Stdin)
		input, _ := stdin.ReadString('\n')
		cmd := strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
		switch cmd[0] {
		case "pause":
			clock <- false
			os.Stdout.WriteString("Simulation paused\033[F")
		case "resume":
			clock <- true
			os.Stdout.WriteString("Simulation resumed\033[F")
		case "rate":
			i, _ := strconv.Atoi(cmd[1])
			if i > 0 && i < 11 {
				for _, inj := range injs {
					inj <- i
				}
				os.Stdout.WriteString("Rate changed to " + cmd[1] + "\033[F")
			} else {
				os.Stdout.WriteString("Rate but be between 1 and 10\033[F")
			}
		case "drunk":
		case "exit":
		default:
			os.Stdout.WriteString("Unknown command\033[F")
		}
	}
}
