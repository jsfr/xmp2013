package xmpmud

import (
	"strconv"
)

func Player(ctrl chan Cmd, ctrlresp chan Cmd, resp chan Cmd, room chan Cmd,
	id string, items map[string]uint64) {
	for ctrl != nil {
		req := (<-ctrl)
		cmd := req.cmd
		switch cmd[0] {
		case "n", "e", "s", "w":
			room <- Cmd{[]string{cmd[0], id}, resp}
			result := <-resp
			switch result.cmd[0] {
			case "locked":
				ctrlresp <- Cmd{[]string{"The door is locked"}, nil}
			case "room":
				room = result.resp
				room <- Cmd{[]string{"look", id}, resp}
				result = <-resp
				invDesc := "Your inventory contains: "
				for i, q := range items {
					invDesc += strconv.FormatUint(q, 10) + " " + i + ", "
				}
				respCmd := "You enter the room and look around\n" +
					result.cmd[0] + "\n" + invDesc
				ctrlresp <- Cmd{[]string{respCmd}, nil}
			}
		case "look":
			room <- Cmd{[]string{"look", id}, resp}
			result := <-resp
			invDesc := "Your inventory contains: "
			for i, q := range items {
				invDesc += strconv.FormatUint(q, 10) + " " + i + ", "
			}
			ctrlresp <- Cmd{[]string{result.cmd[0] + "\n" + invDesc}, nil}
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
			ctrlresp <- Cmd{[]string{text}, nil}
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
			ctrlresp <- Cmd{[]string{text}, nil}
		case "player":
			ctrlresp <- Cmd{[]string{cmd[1] + " entered the room"}, nil}
		case "exit":
			close(ctrlresp)
			ctrl = nil
		default:
			ctrlresp <- Cmd{[]string{"I'm sorry Dave, I'm afraid I can't do that"}, nil}
		}
	}
}
