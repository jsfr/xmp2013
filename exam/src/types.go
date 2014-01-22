package main

type CmdWord int
type Direction int

const (
	CmdPause  CmdWord = iota
	CmdResume CmdWord
	CmdRate   CmdWord
	CmdDrunk  CmdWord
)

const (
	Up    Direction = iota
	Right Direction
	Down  Direction
	Left  Direction
)

type Msg struct {
	Rsp   chan bool
	Agent Pedestrian
}

type Pedestrian struct {
	Id    string
	Dirc  Direction
	Clock chan bool
}

type Coordinate struct {
	X int
	Y int
}
