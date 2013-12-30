package xmpmud

type Cmd struct {
	cmd  []string
	resp chan Cmd
}
