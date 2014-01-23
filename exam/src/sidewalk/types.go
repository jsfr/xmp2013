package sidewalk

type Direction int
type State int
type CmdWord int

const (
	Up Direction = iota
	Right
	Down
	Left
)

const (
	GotDrunk State = iota
	GotSober
	Regular
)

const (
	CmdResume CmdWord = iota
	CmdPause
	CmdTime
	CmdExit
)

type Msg struct {
	Rsp    chan bool
	Agent  Pedestrian
	Status State
}

type Pedestrian struct {
	Dirc Direction
	Ok   bool
}

type Coordinate struct {
	X int
	Y int
}

type DrunkRequest struct {
	ReqCoord  Coordinate
	DrunkDirc Direction
	Rsp       chan bool
}

type ClockRequest struct {
	Cmd CmdWord
	Arg []int
}
