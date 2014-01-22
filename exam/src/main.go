package main

import (
	"os"
	"strconv"
)

func Clock(in chan (chan bool), injs [](chan bool), printer chan bool, steps int) {
	var agents map[string](chan bool)
	for step := 0; step < steps; step++ {
		for c, inj := range injs { // tick every injector
			inj <- true
		}
		for i := 0; i < len(injs); i++ { // as many times as there are injectors check input
			newAgent := <-in
			if newAgent != nil {
				agents = append(agents, newAgent)
			}
		}
		for i, agent := range agents { // tick the clock for agents
			agent <- true
		}
		for i, agent := range agents { // get response from agents
			if !<-agent {
				// The agents informs that it has walked of the board
				delete(agents, i)
			}
		}
		printer <- true // Now everone has moved and the printer can print
	}
}

func Injector(id string, clock chan bool, cell chan Msg, out chan (chan bool), d Direction) {
	step := 3
	n := 0
	for _ = range clock {
		agentClock := nil
		if n%step == 0 { // create a new player
			agentClock = make(chan bool)
			agent := Pedestrian{id + strconv.Itoa(n), d, agentClock}

			rsp := make(chan bool)
			cell <- Msg{rsp, agent}
			if !<-rsp {
				agentClock = nil
			}
		}
		out <- agentClock
		n++
	}
}

func Printer(in chan Coordinate, clock chan bool, x int, y int) {
	coords := makeCoords(x, y)
	step := 0
	for {
		select {
		case c := <-in:
			coords[c.Y][c.X] = true // A player resides at this coordinate in this turn
		case <-clock:
			for i, _ := range coords {
				line := ""
				for j, _ := range coords[i] {
					if coords[i][j] {
						line += "**"
					} else {
						line += "00"
					}
				}
				line += "\n"
				os.Stdout.WriteString("Step " + strconv.Itoa(step) + "\n")
				os.Stdout.WriteString(line)
				os.Stdout.WriteString(line)
			}
			setCoords(false, &coords)
			step++
		}
	}
}

func Cell(in chan Msg, out [](chan Msg), printer chan Coordinate, coord Coordinate) {
	agent := new(Pedestrian)
	oldAgent := false // used to send a print if the agent has been stagnant a turn
	rsp := make(chan bool)
	for {
		if !oldAgent {
			msg := <-in
			agent = msg.Agent
			msg.Rsp <- true
			<-msg.Rsp
			printer <- coord
		}
		select {
		case <-agent.Clock:
			d := directions(agent.Dirc)
			tries := 0
			done := false
			for tries < 3 && !done {
				if out[d[tries]] == nil { // Agent walked off the board
					agent = new(Pedestrian)
					agent.Clock <- false
					oldAgent = false
					done = true
				} else {
					select {
					case out[d[tries]] <- Msg{rsp, agent}:
						if <-rsp {
							agent.Clock <- true
							rsp <- true
							agent = new(Pedestrian)
							oldAgent = false
							done = true
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
				agent.Clock <- true
			}
		case msg := <-in:
			msg.Rsp <- false
		}
	}
}

func Sidewalk(x int, y int, injs []Coordinate, ds []Direction, ids []string, steps int) {
	injChs := make([](chan bool), size(injs))
	clockIn := make(chan (chan bool))
	printerClock := make(chan bool)
	printer := make(chan Coordinate)
	cells := makeMsgChs(x, y)
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			out := getOutChs(j, i, x, y)
			go Cell(cell[i][j], out, printer, Coordinate{j, i})
		}
	}
	for i, inj := range injs {
		injChs[i] = make(chan (bool))
		injx := injs[i].X
		injy := injs[i].Y
		go Injector(ids[i], injChs[i], cell[y][x], clockIn, ds[i])
	}
	go Printer(printer, printerClock, x, y)
	Clock(clockIn, injChs, printerClock, steps)
}

func Main() {
	Sidewalk(5, 10, []Coordinate{Coordinate{3, 0}}, []Direction{Up}, []string{"A"}, 200)
}
