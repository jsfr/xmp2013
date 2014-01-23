package sidewalk

import (
	"math"
	"math/rand"
)

/**
 * Creates a list of the given direction and its two orthogonal directions
 * the orthogonals order are randomized.
 */
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

/**
 * Creates a slice of slices containing bools.
 */
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

/**
 * Sets all values in the slice of slices to the specified value.
 */
func setCoords(val bool, coords *[][]bool) {
	for i, _ := range *coords {
		for j, _ := range (*coords)[i] {
			(*coords)[i][j] = val
		}
	}
}

/**
 * Creates a slice of slices containing message channels used by cells.
 */
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

/**
 * Returns a list of the four channels corresponding to cells adjacent to
 * the cell at the position (j,i). If no cell is adjacent in a direction the
 * channel is nil.
 */
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

/**
 * Checks whether a coordinate is valid (non-negative) and
 * within the specified bounds.
 */
func validCoord(coord Coordinate, x int, y int) bool {
	return coord.X >= 0 && coord.Y >= 0 && coord.X < x && coord.Y < y
}

func randomDirc(coord Coordinate, x int, y int) Coordinate {
	randomCoord := Coordinate{-1, -1}
	switch Direction(rand.Intn(4)) {
	case Up:
		if coord.Y < y-1 {
			randomCoord = Coordinate{coord.X, coord.Y + 1}
		}
	case Down:
		if coord.Y > 0 {
			randomCoord = Coordinate{coord.X, coord.Y - 1}
		}
	case Left:
		if coord.X > 0 {
			randomCoord = Coordinate{coord.X - 1, coord.Y}
		}
	case Right:
		if coord.X < x-1 {
			randomCoord = Coordinate{coord.X + 1, coord.Y}
		}
	}
	return randomCoord
}

/**
 * Calculates whether a given drunk coordinate is within 3 cells.
 * The calculation is the euclidian distance formula in two dimensions.
 */
func drunkNearby(drunk Coordinate, reqCoord Coordinate, x int, y int) bool {
	if validCoord(drunk, x, y) {
		x := float64(drunk.X - reqCoord.X)
		y := float64(drunk.Y - reqCoord.Y)
		dist := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
		return dist <= 3
	} else {
		return false
	}
}
