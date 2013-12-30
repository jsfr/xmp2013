package xmpmud

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Shell(ctrl chan Cmd, resp chan Cmd) {
	stdin := bufio.NewReader(os.Stdin)
	var cmd []string = nil
	for {
		fmt.Print("\n> ")
		input, _ := stdin.ReadString('\n')
		cmd = strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
		ctrl <- Cmd{cmd, resp}
		fmt.Println((<-resp).cmd[0])
	}
}

func Ai(ctrl chan Cmd, resp chan Cmd) {
	ds := []string{"n", "e", "s", "w"}
	for {
		time.Sleep(2 * time.Second)
		d := ds[rand.Int63n(4)]
		ctrl <- Cmd{[]string{d}, resp}
		_ = <-resp
	}
}
