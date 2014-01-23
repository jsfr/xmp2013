package sidewalk

type CmdWord int
type Direction int

const (
	CmdPause CmdWord = iota
	CmdResume
	CmdRate
	CmdDrunk
)

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
	Id    string
	Dirc  Direction
	Clock chan bool
	Done  chan bool
}

type Coordinate struct {
	X int
	Y int
}

type PrintMsg struct {
	Coord Coordinate
	Id    string
}
