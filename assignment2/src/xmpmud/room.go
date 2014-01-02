package xmpmud

import (
	"fmt"
	"strconv"
)

func Room(req chan Cmd, out map[string](chan Cmd), items map[string]uint64,
	players map[string](chan Cmd)) {
	for req != nil {
		c := <-req
		switch c.cmd[0] {
		case "player":
			for _, pchan := range players {
				pchan <- Cmd{[]string{"player", c.cmd[1]}, nil}
			}
			players[c.cmd[1]] = c.resp
		case "look":
			itemsDesc := "The room contains the following items: "
			playersDesc := "The room contains the following players: "
			for i, q := range items {
				itemsDesc += strconv.FormatUint(q, 10) + " " + i + ", "
			}
			for p, _ := range players {
				playersDesc += p + ", "
			}
			c.resp <- Cmd{[]string{itemsDesc + "\n" + playersDesc}, nil}
		case "n", "e", "s", "w":
			if out[c.cmd[0]] != nil {
				out[c.cmd[0]] <- Cmd{[]string{"player", c.cmd[1]},
					players[c.cmd[1]]}
				c.resp <- Cmd{[]string{"room"}, out[c.cmd[0]]}
				delete(players, c.cmd[1])
			} else {
				c.resp <- Cmd{[]string{"locked"}, nil}
			}
		case "take":
			item := c.cmd[2]
			_, ok := items[item]
			if ok {
				items[item] -= 1
				if items[item] == 0 {
					delete(items, item)
				}
				c.resp <- Cmd{[]string{"got"}, nil}
			} else {
				c.resp <- Cmd{[]string{"empty"}, nil}
			}
		case "leave":
			item := c.cmd[2]
			items[item] += 1
			c.resp <- Cmd{nil, nil}
		case "exit":
			req = nil
		default:
			fmt.Println("I always thought something was fundamentally" +
				" wrong with the universe")
			c.resp <- Cmd{nil, nil}
		}
	}
}
