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
	go Room(c[0], map[string](chan Cmd){"n": nil, "e": c[1], "s": c[2], "w": nil},
		map[string]uint64{"gold": 1}, map[string](chan Cmd){"John Doe": c[4]})
	go Room(c[1], map[string](chan Cmd){"n": nil, "e": nil, "s": c[3], "w": c[0]},
		map[string]uint64{"gold": 2}, map[string](chan Cmd){"Slenderman": c[7]})
	go Room(c[2], map[string](chan Cmd){"n": c[0], "e": c[3], "s": nil, "w": nil},
		map[string]uint64{"gold": 3}, map[string](chan Cmd){})
	go Room(c[3], map[string](chan Cmd){"n": c[1], "e": nil, "s": nil, "w": c[2]},
		map[string]uint64{"gold": 4}, map[string](chan Cmd){})
	go Player(c[4], c[5], c[6], c[0], "John Doe", map[string]uint64{})
	go Player(c[7], c[8], c[9], c[1], "Slenderman", map[string]uint64{})
	go Ai(c[7], c[8])
	Shell(c[4], c[5], [](chan Cmd){c[7], c[0], c[1], c[2], c[3]})
}
