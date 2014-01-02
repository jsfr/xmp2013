package xmpmud

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Input(out chan string) {
	stdin := bufio.NewReader(os.Stdin)
	ok := true
	for ok {
		input, _ := stdin.ReadString('\n')
		out <- input
		_, ok = <-out
	}
}

func Shell(ctrl chan Cmd, resp chan Cmd, channels [](chan Cmd)) {
	in := make(chan string)
	go Input(in)
	for resp != nil {
		fmt.Print("\n> ")
		select {
		case input := <-in:
			cmd := strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
			ctrl <- Cmd{cmd, nil}
			msg, ok := <-resp
			if ok {
				fmt.Println(msg.cmd[0])
				in <- ""
			} else {
				resp = nil
			}
		case cmd := <-resp:
			fmt.Println("\n" + cmd.cmd[0])
		}
	}
	for _, c := range channels {
		c <- Cmd{[]string{"exit"}, nil}
	}
	close(in)
}

func Ai(ctrl chan Cmd, resp chan Cmd) {
	ds := []string{"n", "e", "s", "w"}
	ok := true
	for ok {
		time.Sleep(2 * time.Second)
		d := ds[rand.Intn(len(ds))]
		select {
		case <-resp:
		case ctrl <- Cmd{[]string{d}, nil}:
			_, ok = <-resp
		}
	}
}
