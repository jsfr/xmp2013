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

func Cell(in chan Msg, out map[int](chan Msg), coord Coordinate) {
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
	clock chan Tick, tick chan int) {
	in := make(chan Msg)
	exit := false
	for _ = range tick {
		done := false
		tries := 0
		r := rand.Intn(2)
		d := make([]int, 3)
		d[0] = dirc
		d[2-r] = (dirc + 1) % 4
		d[1+r] = (dirc - 1) % 4
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
			fmt.Println("id:", id, "coord:", coord)
			clock <- Tick{id, tick}
		} else {
			clock <- Tick{id, nil}
			close(tick)
		}
	}
}

func Injector(id string, clock chan Tick, tick chan int, dirc int, cell chan Msg) {
	step := 0
	for _ = range tick {
		if step%4 == 0 {
			tick := make(chan int)
			clock <- Tick{id + strconv.Itoa(step/4), tick}
			go Pedestrian(id+strconv.Itoa(step/4), dirc, Coordinate{-1, -1}, cell, clock, tick)
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
}

func Sidewalk() {
	in := make([](chan Msg), 9)
	clock := make(chan Tick)
	clocks := make(map[string](chan int))
	tick0 := make(chan int)
	tick1 := make(chan int)
	clocks["A"] = tick0
	clocks["B"] = tick1
	for i := 0; i < 9; i++ {
		in[i] = make(chan Msg)
	}
	go Cell(in[0], map[int](chan Msg){0: in[3], 1: in[1], 2: nil, 3: nil}, Coordinate{1, 1})
	go Cell(in[1], map[int](chan Msg){0: in[4], 1: in[2], 2: nil, 3: in[0]}, Coordinate{2, 1})
	go Cell(in[2], map[int](chan Msg){0: in[5], 1: nil, 2: nil, 3: in[1]}, Coordinate{3, 1})
	go Cell(in[3], map[int](chan Msg){0: in[6], 1: in[4], 2: in[0], 3: nil}, Coordinate{1, 2})
	go Cell(in[4], map[int](chan Msg){0: in[7], 1: in[5], 2: in[1], 3: in[3]}, Coordinate{2, 2})
	go Cell(in[5], map[int](chan Msg){0: in[8], 1: nil, 2: in[2], 3: in[4]}, Coordinate{3, 2})
	go Cell(in[6], map[int](chan Msg){0: nil, 1: in[7], 2: in[3], 3: nil}, Coordinate{1, 3})
	go Cell(in[7], map[int](chan Msg){0: nil, 1: in[8], 2: in[4], 3: in[6]}, Coordinate{2, 3})
	go Cell(in[8], map[int](chan Msg){0: nil, 1: nil, 2: in[5], 3: in[7]}, Coordinate{3, 3})
	go Injector("A", clock, tick0, 0, in[0])
	go Injector("B", clock, tick1, 0, in[1])
	Clock(clock, clocks, 10)
}

func main() {
	Sidewalk()
}
