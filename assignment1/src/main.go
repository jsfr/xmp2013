package main

import (
	"fmt"
)

type binMsg struct {
	idx   int
	balls int
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
	layers = 10
	msgList = make([]int, layers+1)
	go nLayerMachineOsc(balls, layers, 0.0, 0.5, 0.1, messages)
	for i := 0; i < layers+1; i++ {
		msg := <-messages
		msgList[msg.idx] = msg.balls
	}
	fmt.Println(layers, "Layer Machine, Oscillating:")
	for i, b := range msgList {
		fmt.Println("Bin:", i, " Balls:", b)
	}

}
