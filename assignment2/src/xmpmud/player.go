package xmpmud

import (
	"fmt"
	"strconv"
)

func Player(resp chan Cmd, room chan Cmd, id string, items map[string]uint64, ctrl chan Cmd) {
	for {
		req := (<-ctrl)
		cmd := req.cmd
		switch cmd[0] {
		case "n", "e", "s", "w":
			room <- Cmd{[]string{cmd[0], id}, resp}
			result := <-resp
			switch result.cmd[0] {
			case "locked":
				req.resp <- Cmd{[]string{"The door is locked"}, nil}
			case "room":
				room = result.resp
				room <- Cmd{[]string{"look", id}, resp}
				result = <-resp
				invDesc := "Your inventory contains: "
				for i, q := range items {
					invDesc += strconv.FormatUint(q, 10) + " " + i + ", "
				}
				req.resp <- Cmd{[]string{"You enter the room and look around\n" + result.cmd[0] + "\n" + invDesc}, nil}
			}
		case "look":
			room <- Cmd{[]string{"look", id}, resp}
			result := <-resp
			invDesc := "Your inventory contains: "
			for i, q := range items {
				invDesc += strconv.FormatUint(q, 10) + " " + i + ", "
			}
			req.resp <- Cmd{[]string{result.cmd[0] + "\n" + invDesc}, nil}
		case "take":
			item := cmd[1]
			room <- Cmd{[]string{"take", id, item}, resp}
			result := <-resp
			text := ""
			switch result.cmd[0] {
			case "got":
				_, ok := items[item]
				if ok {
					items[item] += 1
				} else {
					items[item] = 1
				}
				text = "You took 1 " + item + " from the room"
			case "empty":
				text = "There is no " + item + " in the room"
			}
			req.resp <- Cmd{[]string{text}, nil}
		case "leave":
			item := cmd[1]
			_, ok := items[item]
			text := ""
			if ok {
				room <- Cmd{[]string{"leave", id, item}, resp}
				_ = <-resp
				items[item] -= 1
				if items[item] == 0 {
					delete(items, item)
				}
				text = "You left 1 " + item + " in the room"
			} else {
				text = "Your inventory does not contain any " + item
			}
			req.resp <- Cmd{[]string{text}, nil}
		case "exit":
			fallthrough // TODO
		case "player":
			fmt.Println("\n", cmd[1], "entered the same room as", id, "\n> ")
		default:
			req.resp <- Cmd{[]string{"I'm sorry Dave, I'm afraid I can't do that"}, nil}
		}
	}
}
