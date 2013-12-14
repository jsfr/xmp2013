package main

func oscillator(in chan int, outs [](chan float64)) {
	delta := 0.1
	min := 0.0
	max := 0.5
	val := 0.0
	sign := 1.0
	c := true
	for c {
		if val >= max {
			sign = -1.0
		}
		if val <= min {
			sign = 1.0
		}
		val += sign * delta
		_, c := <-in

		for _, ch := range outs {
			if c {
				select {
				case ch <- val:
				default:
				}
			} else {
				close(ch)
			}
		}
	}
}

func injector(out chan int, balls int) {
	for i := 0; i < balls; i++ {
		out <- 1
	}
	close(out)
}
