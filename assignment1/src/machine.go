package main

import (
	"fmt"
	"math/rand"
)

type binMsg struct {
	idx   int
	balls int
}

// TODO: Add shutdown to all processes

func injector(out chan int, balls int) {
	for i := 0; i < balls; i++ {
		out <- 1
	}
	close(out)
}

func pin(in0 chan int, in1 chan int, out0 chan int, out1 chan int) {
	var choice int
	for in0 != nil || in1 != nil {
		choice = rand.Intn(2)
		select {
		case _, c0 := <-in0:
			if !c0 {
				in0 = nil
			} else {
				if choice == 0 {
					out0 <- 1
				} else {
					out1 <- 1
				}
			}
		case _, c1 := <-in1:
			if !c1 {
				in1 = nil
			} else {
				if choice == 0 {
					out0 <- 1
				} else {
					out1 <- 1
				}
			}
		}
	}
	if out0 != out1 {
		close(out0)
	}
	close(out1)
}

func oscilator(in chan int, outs [](chan int)) {
	delta := 0.1
	min := 0
	max := 0.5
	val := 0
	sign := 1
	for {
		if val == max {
			sign = -1
		}
		if val == min {
			sign = 1
		}
		val += sign * delta
		_, c := <-in

		if c {
			for _, ch := range outs {
				if c {
					ch <- val
				} else {
					close(ch)
				}
			}
		}
	}
}

func pinOsc(in0 chan int, in1 chan int, out0 chan int, out1 chan int, osc chan int) {
	var choice int
	for in0 != nil || in1 != nil {
		choice = rand.Intn(2)
		select {
		case _, c0 := <-in0:
			if !c0 {
				in0 = nil
			} else {
				if choice == 0 {
					out0 <- 1
				} else {
					out1 <- 1
				}
			}
		case _, c1 := <-in1:
			if !c1 {
				in1 = nil
			} else {
				if choice == 0 {
					out0 <- 1
				} else {
					out1 <- 1
				}
			}
		}
	}
	if out0 != out1 {
		close(out0)
	}
	close(out1)
}

func bin(in0 chan int, in1 chan int, n int, message chan binMsg) {
	var balls int = 0
	for in0 != nil || in1 != nil {
		select {
		case _, c0 := <-in0:
			if !c0 {
				in0 = nil
			} else {
				balls++
			}
		case _, c1 := <-in1:
			if !c1 {
				in1 = nil
			} else {
				balls++
			}
		}
	}
	message <- binMsg{n, balls}
}

func twoLayerMachine(balls int, messages chan binMsg) {
	n := 10
	chs := make([](chan int), n)
	for i := 0; i < n; i++ {
		chs[i] = make(chan int)
	}

	go pin(chs[0], nil, chs[1], chs[2])
	go pin(chs[1], chs[1], chs[3], chs[4])
	go pin(chs[2], chs[2], chs[5], chs[6])
	go pin(chs[3], chs[3], chs[7], chs[7])
	go pin(chs[4], chs[5], chs[8], chs[8])
	go pin(chs[6], chs[6], chs[9], chs[9])
	go bin(chs[7], nil, 0, messages)
	go bin(chs[8], nil, 1, messages)
	go bin(chs[9], nil, 2, messages)
	go injector(chs[0], balls)
}

func nLayerMachine(balls int, layers int, messages chan binMsg) {
	chs := make([][](chan int), layers+2)
	chs[0] = make([](chan int), 2)
	chs[0][0] = make(chan int)
	chs[0][1] = nil

	for i := 1; i < layers+2; i++ {
		chs[i] = make([](chan int), 2*i+2)
		chs[i][0] = nil
		chs[i][2*i+1] = nil
		for j := 1; j < 2*i+1; j++ {
			chs[i][j] = make(chan int)
		}
		for j := 0; j < i; j++ {
			go pin(chs[i-1][2*j], chs[i-1][2*j+1], chs[i][2*j+1], chs[i][2*j+2])
		}
	}
	for i := 0; i < layers+1; i++ {
		go bin(chs[layers+1][2*i+1], chs[layers+1][2*i+2], i, messages)
	}
	go injector(chs[0][0], balls)
}

func main() {
	var messages chan binMsg = make(chan binMsg)
	var msgList []int
	var layers int
	var balls int = 1000

	// Problem 1: Fixed Two-Layer Machine Simulation
	msgList = make([]int, 3)
	go twoLayerMachine(balls, messages)
	for n := 0; n < 3; n++ {
		msg := <-messages
		msgList[msg.idx] = msg.balls
	}
	fmt.Println("Fixed Two-Layer Machine:")
	for i, b := range msgList {
		fmt.Println("Bin:", i, "Balls:", b)
	}

	// Problem 2: N-Layer Machine Simulation (set to 10 layers)
	layers = 10
	msgList = make([]int, layers+1)
	go nLayerMachine(balls, layers, messages)
	for i := 0; i < layers+1; i++ {
		msg := <-messages
		msgList[msg.idx] = msg.balls
	}
	fmt.Println(layers, "Layer Machine:")
	for i, b := range msgList {
		fmt.Println("Bin:", i, " Balls:", b)
	}

	// Problem 3: Oscilator
}
