package main

import (
	"fmt"
	"math/rand"
	"strconv"
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
			}
		case _, c1 := <-in1:
			if !c1 {
				in1 = nil
			}
		}
		if choice == 0 {
			out0 <- 1
		} else {
			out1 <- 1
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
			}
		case _, c1 := <-in1:
			if !c1 {
				in1 = nil
			}
		}
		balls++
	}
	message <- binMsg{n, balls}
}

func twoLayerMachine(balls int, messages chan binMsg) {
	n := 10
	chs := make([](chan int), n)
	for i := 0; i < n; i++ {
		chs[i] = make(chan int)
	}

	go pin(chs[0], chs[0], chs[1], chs[2])
	go pin(chs[1], chs[1], chs[3], chs[4])
	go pin(chs[2], chs[2], chs[5], chs[6])
	go pin(chs[3], chs[3], chs[7], chs[7])
	go pin(chs[4], chs[5], chs[8], chs[8])
	go pin(chs[6], chs[6], chs[9], chs[9])
	go bin(chs[7], chs[7], 0, messages)
	go bin(chs[8], chs[8], 1, messages)
	go bin(chs[9], chs[9], 2, messages)
	go injector(chs[0], balls)
}

func nLayerMachine(balls int, layers int, messages chan binMsg) {
	chs := make([][](chan int), layers+2)
	chs[0] = make([](chan int), 2)
	chs[0][0] = make(chan int)
	chs[0][1] = chs[0][0]

	for i := 0; i < layers+1; i++ {
		chs[i+1] = make([](chan int), 2*i+4)
		chs[i+1][0] = nil
		chs[i+1][2*i+3] = nil
		for j := 1; j < 2*i+3; j++ {
			chs[i+1][j] = make(chan int)
		}
		for j := 0; j < i+1; j++ {
			go pin(chs[i][2*j], chs[i][2*j+1], chs[i+1][2*j+1], chs[i+1][2*j+2])
		}
	}
	for i := 0; i < layers+1; i++ {
		go bin(chs[layers+1][2*i+1], chs[layers+1][2*i+2], i, messages)
	}
	go injector(chs[0][0], balls)
}

func main() {
	/*	go twoLayerMachine(1000, messages)
		for n := 0; n < 3; n++ {
			msg := <-messages
			fmt.Println(msg)
		}*/

	n, iter := 10, 1000
	messages := make(chan binMsg)
	msgList := make([]int, n+1)
	go nLayerMachine(iter, n, messages)
	for i := 0; i < n+1; i++ {
		msg := <-messages
		msgList[msg.idx] = msg.balls
	}
	fmt.Println(n, " Layer Machine:")
	for i, b := range msgList {
		fmt.Println("Bin: ", strconv.Itoa(i), " Balls:", strconv.Itoa(b))
	}
}
