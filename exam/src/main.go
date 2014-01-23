package main

import (
	. "./sidewalk"
)

func main() {
	Sidewalk(5, 10, []Coordinate{Coordinate{2, 0}, Coordinate{2, 9}}, []Direction{Up, Down}, 10)
}
