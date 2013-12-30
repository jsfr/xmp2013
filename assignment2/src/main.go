package main

import (
	. "./xmpmud"
)

func main() {
	n := 10
	c := make([](chan Cmd), n)
	for i := 0; i < n; i++ {
		c[i] = make(chan Cmd)
	}
	go Room(c[1], map[string](chan Cmd){"n": nil, "e": c[2], "s": c[3], "w": nil}, map[string]uint64{"gold": 1}, map[string](chan Cmd){"John Doe": c[5]})
	go Room(c[2], map[string](chan Cmd){"n": nil, "e": nil, "s": c[4], "w": c[1]}, map[string]uint64{"gold": 2}, map[string](chan Cmd){"Slenderman": c[8]})
	go Room(c[3], map[string](chan Cmd){"n": c[1], "e": c[4], "s": nil, "w": nil}, map[string]uint64{"gold": 3}, map[string](chan Cmd){})
	go Room(c[4], map[string](chan Cmd){"n": c[2], "e": nil, "s": nil, "w": c[3]}, map[string]uint64{"gold": 4}, map[string](chan Cmd){})
	go Player(c[0], c[1], "John Doe", map[string]uint64{}, c[5])
	go Player(c[7], c[2], "Slenderman", map[string]uint64{}, c[8])
	go Ai(c[8], c[9])
	Shell(c[5], c[6])
}
