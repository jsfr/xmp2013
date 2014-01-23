package sidewalk

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Msg struct {
	Rsp   chan bool
	Agent Pedestrian
}

type Pedestrian struct {
	Id   string
	Dirc Direction
}

type Coordinate struct {
	X int
	Y int
}

type PrintMsg struct {
	Coord Coordinate
	Id    string
}
