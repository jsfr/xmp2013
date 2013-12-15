package main

import (
	"math/rand"
)

// injector runs through the number of balls,
// sending a message to the out channel every time.
// Upon finish the out channel is closed
func injector(out chan int, balls int) {
	for i := 0; i < balls; i++ {
		out <- 1
	}
	close(out)
}

// pin receives balls on its in channels, and depending on the random choice
// passes it on to its respective out channel. If an in channel closes
// it is nil'ed to keep it from always being ready. The pin closes its out
// channels when both in channels have been nil'ed (i.e. closed)
func pin(in0 chan int, in1 chan int, out0 chan int, out1 chan int) {
	for in0 != nil || in1 != nil {
		choice := rand.Intn(2)
		c := false
		select {
		case _, c0 := <-in0:
			if c0 {
				c = true
			} else {
				in0 = nil
			}
		case _, c1 := <-in1:
			if c1 {
				c = true
			} else {
				in1 = nil
			}
		}
		if c && choice <= 0 {
			out0 <- 1
		} else if c && choice > 0 {
			out1 <- 1
		}
	}
	if out0 != out1 {
		close(out0)
	}
	close(out1)
}

// bin receives balls on its in channels and counts how many it is currently
// holding. When a channel is closed it is nil'ed to keep it from always being
// ready. When both channels are closed a message is sent with the bin count on
// the message channel
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

// A static triangular shaped two-layered machine
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

// A dynamic triangular n-layered machine
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

// oscillator waits on request and ball channel. The ball channel updates
// the current skew value, whereas the request channel is used to request
// the current skew.
func oscillator(ball chan int, request chan (chan float64), min float64,
	max float64, delta float64) {
	c0 := true
	val := 0.0

	for c0 {
		select {
		case _, c := <-ball:
			switch {
			case c && val >= max:
				val -= delta
			case c && val <= min:
				val += delta
			case !c:
				//TODO: How to close this?
			}
		case resp, c := <-request:
			switch {
			case c:
				resp <- val
			case !c:
				//TODO: How to close this?
			}
		}
	}
}

// injectorOsc works as injector, but it also sends a message on a oscillator
// channel to update the current skew. This is done first, to make sure the skew
// is updated when the first pin requests it from the oscillator.
func injectorOsc(out chan int, osc chan int, balls int) {
	for i := 0; i < balls; i++ {
		osc <- 1
		out <- 1
	}
	//close(osc)
	close(out)
}

// pinOsc works as pin, but before the ball is sent further down the pin
// requests the current skew from the oscillator by sending its reponse channel
// on the request channel.
func pinOsc(in0 chan int, in1 chan int, out0 chan int, out1 chan int,
	osc chan (chan float64)) {
	resp := make(chan float64)
	for in0 != nil || in1 != nil {
		choice := 2 * rand.Float64()
		c := false
		select {
		case _, c0 := <-in0:
			if c0 {
				c = true
			} else {
				in0 = nil
			}
		case _, c1 := <-in1:
			if c1 {
				c = true
			} else {
				in1 = nil
			}
		}
		if c {
			osc <- resp
			skew, c0 := <-resp
			if c0 {
				choice += skew
			}
			if choice <= 1 {
				out0 <- 1
			} else if choice > 1 {
				out1 <- 1
			}
		}
	}
	if out0 != out1 {
		close(out0)
	}
	close(out1)
}

// a dynamic n-layered oscillating machine
func nLayerMachineOsc(balls int, layers int, min float64, max float64,
	delta float64, messages chan binMsg) {
	chs := make([][](chan int), layers+2)
	chs[0] = make([](chan int), 2)
	chs[0][0] = make(chan int)
	chs[0][1] = nil
	injosc := make(chan int)
	pinosc := make(chan (chan float64))

	for i := 1; i < layers+2; i++ {
		chs[i] = make([](chan int), 2*i+2)
		chs[i][0] = nil
		chs[i][2*i+1] = nil
		for j := 1; j < 2*i+1; j++ {
			chs[i][j] = make(chan int)
		}
		for j := 0; j < i; j++ {
			go pinOsc(chs[i-1][2*j], chs[i-1][2*j+1], chs[i][2*j+1],
				chs[i][2*j+2], pinosc)
		}
	}
	for i := 0; i < layers+1; i++ {
		go bin(chs[layers+1][2*i+1], chs[layers+1][2*i+2], i, messages)
	}
	go oscillator(injosc, pinosc, min, max, delta)
	go injectorOsc(chs[0][0], injosc, balls)
}
