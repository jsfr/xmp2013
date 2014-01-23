package sidewalk

import (
	"math/rand"
)

func directions(d Direction) (l []Direction) {
	r := rand.Intn(2)
	switch d {
	case Up:
		if r > 0 {
			l = []Direction{Up, Left, Right}
		} else {
			l = []Direction{Up, Right, Left}
		}
	case Down:
		if r > 0 {
			l = []Direction{Down, Right, Left}
		} else {
			l = []Direction{Down, Right, Left}
		}
	case Left:
		if r > 0 {
			l = []Direction{Left, Down, Up}
		} else {
			l = []Direction{Left, Up, Down}
		}
	case Right:
		if r > 0 {
			l = []Direction{Right, Down, Up}
		} else {
			l = []Direction{Right, Up, Down}
		}
	}
	return
}

func makeCoords(x int, y int) [][]string {
	coords := make([][]string, y)
	for i, _ := range coords {
		coords[i] = make([]string, x)
		for j, _ := range coords[i] {
			coords[i][j] = ""
		}
	}
	return coords
}

func setCoords(val string, coords *[][]string) {
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
	var left, right, up, down chan Msg = nil, nil, nil, nil
	if j == 0 {
		right = chs[i][j+1]
	} else if j == x-1 {
		left = chs[i][j-1]
	} else {
		left = chs[i][j-1]
		right = chs[i][j+1]
	}
	if i == 0 {
		up = chs[i+1][j]
	} else if i == y-1 {
		down = chs[i-1][j]
	} else {
		down = chs[i-1][j]
		up = chs[i+1][j]
	}
	return [](chan Msg){up, right, down, left}
}
