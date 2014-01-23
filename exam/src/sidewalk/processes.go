package sidewalk

import (
	"os"
	"strconv"
	"time"
)

func Clock(in chan Pedestrian, agentClock chan bool, agentDone chan bool,
	injs [](chan bool), printer chan bool, steps int, controller chan bool) {
	agents := 0
	<-controller
	for step := 0; step <= steps; step++ {
		select {
		case play := <-controller:
			if !play {
				for !play {
					play = <-controller // wait to be unpaused
				}
			}
		default:
		}
		for i := 0; i < agents; i++ { // tick the clock for agents
			agentClock <- true
		}
		agentRemove := 0
		for i := 0; i < agents; i++ { // get response from agents
			if !<-agentDone {
				agentRemove++
			}
		}
		agents -= agentRemove
		for _, inj := range injs { // tick every injector
			inj <- true
		}
		for i := 0; i < len(injs); i++ { // as many times as there are injectors check input
			agentCh := <-in
			if agentCh.Id != "" {
				agents++
			}
		}
		time.Sleep(500 * time.Millisecond)
		printer <- true // Now everone has moved and the printer can print
		<-printer
	}
}

func Injector(id string, clock chan bool, cell chan Msg, out chan Pedestrian, d Direction, ctrl chan int) {
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
				agent = Pedestrian{id + strconv.Itoa(n), d}
				cell <- Msg{rsp, agent}
				if !<-rsp {
					agent.Id = ""
				}
			}
			out <- agent
			n++
		}
	}
}

func Printer(in chan PrintMsg, clock chan bool, x int, y int) {
	coords := makeCoords(x, y)
	step := 0
	for {
		select {
		case c := <-in:
			coords[c.Coord.Y][c.Coord.X] = c.Id // A player resides at this coordinate in this turn
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
					if coords[len(coords)-i-1][j] != "" {
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
			setCoords("", &coords)
			step++
			clock <- true
		}
	}
}

func Cell(in chan Msg, out [](chan Msg), printer chan PrintMsg, coord Coordinate,
	agentClock chan bool, agentDone chan bool) {
	var agent Pedestrian
	oldAgent := false // used to send a print if the agent has been stagnant a turn
	rsp := make(chan bool)
	for {
		if !oldAgent {
			msg := <-in
			agent = msg.Agent
			msg.Rsp <- true
			printer <- PrintMsg{coord, agent.Id}
		}
		done := false
		for !done {
			select {
			case <-agentClock:
				d := directions(agent.Dirc)
				tries := 0

				for tries < 3 && !done {
					if out[d[tries]] == nil { // Agent walked off the board
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
						printer <- PrintMsg{coord, agent.Id}
					}
					agentDone <- true
				}
			case msg := <-in:
				msg.Rsp <- false
			}
		}
	}
}

func Sidewalk(x int, y int, injs []Coordinate, ds []Direction, steps int) {
	clockCtrl := make(chan bool)
	injChs := make([](chan bool), len(injs))
	injCtrl := make([](chan int), len(injs))
	clockIn := make(chan Pedestrian)
	printerClock := make(chan bool)
	agentClock := make(chan bool)
	agentDone := make(chan bool)
	printer := make(chan PrintMsg)
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
		go Injector(strconv.Itoa(i), injChs[i], cells[injy][injx], clockIn, ds[i], injCtrl[i])
	}
	go Printer(printer, printerClock, x, y)
	go Clock(clockIn, agentClock, agentDone, injChs, printerClock, steps, clockCtrl)
	Controller(clockCtrl, injCtrl, y)
}
