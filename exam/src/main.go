package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Msg struct {
	Id        string
	Direction int
	Rsp       chan Msg
	Succ      bool
	Coord     Coordinate
}

type Coordinate struct {
	X int
	Y int
}

type Tick struct {
	Id   string
	Tick chan int
}

type Print struct {
	Id     string
	Coords chan Coordinate
}

func Cell(in chan Msg, out [](chan Msg), coord Coordinate) {
	pedestrian := false
	dirc := -1
	newMsg := false
	msg := Msg{"", 0, nil, false, coord}
	for {
		if dirc >= 0 {
			select {
			case m := <-in:
				msg = m
				newMsg = true
			case out[msg.Direction] <- msg:
				dirc = -1
			}
		} else {
			msg = <-in
			newMsg = true
		}
		if newMsg {
			switch {
			case msg.Coord == coord && pedestrian && msg.Rsp != nil:
				// Pedestrian wishes to move
				if out[msg.Direction] != nil {
					dirc = msg.Direction
				} else {
					msg.Rsp <- Msg{"", 0, nil, true, coord}
				}

			case msg.Coord == coord && pedestrian && msg.Rsp == nil:
				// Pedestrian has moved
				pedestrian = false

			case msg.Coord != coord && pedestrian:
				// Already occupied, pedestrian rejected
				msg.Rsp <- Msg{"", 0, in, false, coord}

			case msg.Coord != coord && !pedestrian:
				// Free, pedestrian accepted
				msg.Rsp <- Msg{"", 0, in, true, coord}
				pedestrian = true
			}
			newMsg = false
		}
	}
}

func Pedestrian(id string, dirc int, coord Coordinate, cell chan Msg,
	clock chan Tick, tick chan int, printer chan Coordinate) {
	in := make(chan Msg)
	exit := false
	for _ = range tick {
		done := false
		tries := 0
		r := rand.Intn(2)
		d := make([]int, 3)
		d[0] = dirc
		d[2-r] = (dirc + 3) % 4
		d[1+r] = (dirc + 5) % 4
		for !done && tries <= 2 {
			cell <- Msg{id, d[tries], in, true, coord}
			msg := <-in
			switch {
			case msg.Succ && msg.Rsp != nil:
				if coord.X >= 0 && coord.Y >= 0 {
					cell <- Msg{id, d[tries], nil, true, coord}
				}
				cell = msg.Rsp
				coord = msg.Coord
				done = true
			case msg.Succ && msg.Rsp == nil:
				done = true
				exit = true
			case !msg.Succ:
				tries++
			}
		}
		if !exit {
			printer <- coord
			clock <- Tick{id, tick}
		} else {
			clock <- Tick{id, nil}
			close(tick)
		}
	}
	close(printer)
}

func Printer(id string, in chan Print, clock chan Tick, tick chan int) {
	pedestrians := make(map[string](chan Coordinate))
	k := 0
	for {
		select {
		case <-tick:
			fmt.Println("######", "step", k, "######")
			for i, c := range pedestrians {
				coord, ok := <-c
				if ok {
					fmt.Println("id:", i, "coords:", coord)
				} else {
					fmt.Println("id", i, "walked of board")
					delete(pedestrians, i)
				}
			}
			clock <- Tick{id, tick}
			k++
		case p := <-in:
			pedestrians[p.Id] = p.Coords
		}
	}
}

func Injector(id string, clock chan Tick, tick chan int, dirc int, cell chan Msg,
	stdout chan Print) {
	step := 0
	for _ = range tick {
		if step%4 == 0 {
			tick := make(chan int)
			printer := make(chan Coordinate)
			clock <- Tick{id + strconv.Itoa(step/4), tick}
			stdout <- Print{id + strconv.Itoa(step/4), printer}
			go Pedestrian(id+strconv.Itoa(step/4), dirc, Coordinate{-1, -1}, cell, clock, tick, printer)
		}
		clock <- Tick{id, tick}
		step++
	}
}

func Clock(in chan Tick, clocks map[string](chan int), steps int) {
	for n := 0; n < steps; n++ {
		for _, c := range clocks {
			c <- 1
		}
		k := 0
		for k < len(clocks) {
			tick := <-in
			_, ok := clocks[tick.Id]
			switch {
			case !ok && tick.Tick != nil:
				clocks[tick.Id] = tick.Tick
				k++
			case ok && tick.Tick != nil:
				k++
			case ok && tick.Tick == nil:
				delete(clocks, tick.Id)
			}
		}
	}
	for _, c := range clocks {
		close(c)
	}
}

func Sidewalk() {
	x := 5
	y := 10
	cell := make([]([](chan Msg)), x+2)
	cell[0] = make([](chan Msg), y+2)
	cell[x+1] = make([]chan Msg, y+2)
	for i := 1; i < x+1; i++ {
		cell[i] = make([](chan Msg), y+2)
		for j := 1; j < y+1; j++ {
			cell[i][j] = make(chan Msg)
		}
	}
	for i := 1; i < x+1; i++ {
		for j := 1; j < y+1; j++ {
			out := [](chan Msg){cell[i][j+1], cell[i+1][j], cell[i][j-1], cell[i-1][j]}
			go Cell(cell[i][j], out, Coordinate{i, j})
		}
	}
	tick := make([](chan int), 3)
	for i := 0; i < 3; i++ {
		tick[i] = make(chan int)
	}
	clock := make(chan Tick)
	clocks := make(map[string](chan int))
	stdout := make(chan Print)
	clocks["A"] = tick[0]
	clocks["B"] = tick[1]
	clocks["Printer"] = tick[2]
	go Printer("Printer", stdout, clock, tick[2])
	go Injector("A", clock, tick[0], 0, cell[4][1], stdout)
	go Injector("B", clock, tick[1], 1, cell[1][4], stdout)
	Clock(clock, clocks, 20)
}

func main() {
	Sidewalk()
}
