package main

import (
	"math/rand"
)

func directions(d Direction) []Direction {
	r := rand.Intn(2)
	switch d {
	case Up:
		if r > 0 {
			return []Direction{Up, Left, Right}
		} else {
			return []Direction{Up, Right, Left}
		}
	case Down:
		if r > 0 {
			return []Direction{Down, Right, Left}
		} else {
			return []Direction{Down, Right, Left}
		}
	case Left:
		if r > 0 {
			return []Direction{Left, Down, Up}
		} else {
			return []Direction{Left, Up, Down}
		}
	case Right:
		if r > 0 {
			return []Direction{Right, Down, Up}
		} else {
			return []Direction{Right, Up, Down}
		}
	}
}

func makeCoords(x int, y int) [][]bool {
	coords := make([][]bool, y)
	for i, _ := range coords {
		coords[i] = make([]bool, x)
		for j, _ := range coords[i] {
			coords[i][j] = false
		}
	}
	return coords
}

func setCoords(val bool, coords *[][]bool) {
	for i, _ := range *coords {
		for j, _ := range (*coords)[i] {
			(*coords)[i][j] = val
		}
	}
}

func makeMsgChs(x int, y int) [][](chan Msg) {
	chs := make([][](chan Msg), y)
	for i, _ := range chs {
		chs[i] = make([](chan Msg), x)
		for j, _ := range chs[i] {
			chs[i][j] = make(chan Msg)
		}
	}
	return chs
}

func getOutChs(j int, i int, x int, y int, chs [][](chan Msg)) [](chan Msg) {
	if j == 0 {
		left := nil
		right := chs[i][j+1]
	} else if j == y-1 {
		left := chs[i][j-1]
		right := nil
	} else {
		left := chs[i][j-1]
		right := chs[i][j+1]
	}
	if i == 0 {
		down := nil
		up := chs[i+1][j]
	} else if j == x-1 {
		down := chs[i-1][j]
		up := nil
	} else {
		down := chs[i-1][j]
		up := chs[i+1][j]
	}
	return [](chan Msg){up, left, down, right}
}
