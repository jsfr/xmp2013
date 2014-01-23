package sidewalk

import (
	"container/list"
	"os"
	"strconv"
	"time"
)

func Clock(in chan Pedestrian, injs [](chan bool), printer chan bool, steps int) {
	agents := list.New()
	for step := 0; step <= steps; step++ {
		for e := agents.Front(); e != nil; e = e.Next() { // tick the clock for agents
			e.Value.(Pedestrian).Clock <- true
		}
		for e := agents.Front(); e != nil; { // get response from agents
			tmp := e.Next()
			if !<-e.Value.(Pedestrian).Done {
				// The agents informs that it has walked of the board
				agents.Remove(e)
			}
			e = tmp
		}
		for _, inj := range injs { // tick every injector
			inj <- true
		}
		for i := 0; i < len(injs); i++ { // as many times as there are injectors check input
			agentCh := <-in
			if agentCh.Id != "" {
				agents.PushBack(agentCh)
			}
		}
		time.Sleep(1 * time.Second)
		printer <- true // Now everone has moved and the printer can print
		<-printer
	}
}

func Injector(id string, clock chan bool, cell chan Msg, out chan Pedestrian, d Direction) {
	step := 3
	n := 0
	rsp := make(chan bool)
	for _ = range clock {
		var agent Pedestrian
		if n%step == 0 { // create a new player
			agent = Pedestrian{id + strconv.Itoa(n), d, make(chan bool), make(chan bool)}
			cell <- Msg{rsp, agent}
			if !<-rsp {
				agent.Id = ""
			}
		}
		out <- agent
		n++
	}
}

func Printer(in chan PrintMsg, clock chan bool, x int, y int) {
	coords := makeCoords(x, y)
	step := 0
	os.Stdout.WriteString("\033[0;0H\033[J\033[" + strconv.Itoa(y+2) + ";0H")
	for {
		select {
		case c := <-in:
			coords[c.Coord.Y][c.Coord.X] = c.Id // A player resides at this coordinate in this turn
		case <-clock:
			os.Stdout.WriteString("\033[s\033[0;0HStep " + strconv.Itoa(step) + "\n")
			for i, _ := range coords {
				line := ""
				for j, _ := range coords[len(coords)-i-1] {
					if coords[len(coords)-i-1][j] != "" {
						line += "*"
					} else {
						line += "0"
					}
				}
				line += "\n"
				os.Stdout.WriteString(line)

			}
			os.Stdout.WriteString("\033[u")
			setCoords("", &coords)
			step++
			clock <- true
		}
	}
}

func Cell(in chan Msg, out [](chan Msg), printer chan PrintMsg, coord Coordinate) {
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
			case <-agent.Clock:
				d := directions(agent.Dirc)
				tries := 0

				for tries < 3 && !done {
					if out[d[tries]] == nil { // Agent walked off the board
						agent.Done <- false
						oldAgent = false
						done = true
					} else {
						select {
						case out[d[tries]] <- Msg{rsp, agent}:
							if <-rsp {
								oldAgent = false
								done = true
								agent.Done <- true
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
					agent.Done <- true
				}
			case msg := <-in:
				msg.Rsp <- false
			}
		}
	}
}

func Sidewalk(x int, y int, injs []Coordinate, ds []Direction, steps int) {
	injChs := make([](chan bool), len(injs))
	clockIn := make(chan Pedestrian)
	printerClock := make(chan bool)
	printer := make(chan PrintMsg)
	cells := makeMsgChs(x, y)
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			out := getOutChs(j, i, x, y, cells)
			go Cell(cells[i][j], out, printer, Coordinate{j, i})
		}
	}
	for i, _ := range injs {
		injChs[i] = make(chan (bool))
		injx := injs[i].X
		injy := injs[i].Y
		go Injector(strconv.Itoa(i), injChs[i], cells[injy][injx], clockIn, ds[i])
	}
	go Printer(printer, printerClock, x, y)
	Clock(clockIn, injChs, printerClock, steps)
}
