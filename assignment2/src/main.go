package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Doors struct {
	n, e, s, w chan Cmd
}

type Cmd struct {
	cmd  []string
	resp chan Cmd
}

func room(req chan Cmd, out Doors, items map[string]uint64,
	players map[string](chan Cmd)) {
	for {
		c := <-req
		switch c.cmd[0] {
		case "player":
			// TODO inform other players
			players[c.cmd[1]] = c.resp
			c.resp <- Cmd{[]string{"room"}, req}
		case "look":
			resp := players[c.cmd[1]]
			itemsDesc := "The room contains the following items: "
			playersDesc := "The room contains the following players: "
			for i, q := range items {
				itemsDesc += strconv.FormatUint(q, 10) + " " + i + ", "
			}
			for p, _ := range players {
				playersDesc += p + ", "
			}
			resp <- Cmd{[]string{itemsDesc + "\n" + playersDesc}, nil}
		case "n":
			resp := players[c.cmd[1]]
			if out.n != nil {
				out.n <- Cmd{[]string{"player", c.cmd[1]}, resp}
				delete(players, c.cmd[1])
			} else {
				resp <- Cmd{[]string{"locked"}, nil}
			}
		case "e":
			resp := players[c.cmd[1]]
			if out.e != nil {
				out.e <- Cmd{[]string{"player", c.cmd[1]}, resp}
				delete(players, c.cmd[1])
			} else {
				resp <- Cmd{[]string{"locked"}, nil}
			}
		case "s":
			resp := players[c.cmd[1]]
			if out.s != nil {
				out.s <- Cmd{[]string{"player", c.cmd[1]}, resp}
				delete(players, c.cmd[1])
			} else {
				resp <- Cmd{[]string{"locked"}, nil}
			}
		case "w":
			resp := players[c.cmd[1]]
			if out.w != nil {
				out.w <- Cmd{[]string{"player", c.cmd[1]}, resp}
				delete(players, c.cmd[1])
			} else {
				resp <- Cmd{[]string{"locked"}, nil}
			}
		case "take":
			resp := players[c.cmd[1]]
			item := c.cmd[2]
			quantity, ok := items[item]
			if ok {
				quantity -= 1
				if quantity > 1 {
					items[item] = quantity
				} else {
					delete(items, item)
				}
			}
			/*
				case "leave":
				case "exit":
			*/
		}
	}
}

func player(resp chan Cmd, room chan Cmd, id string, items map[string]uint64) {
	stdin := bufio.NewReader(os.Stdin)
	for {
		input, _ := stdin.ReadString('\n')
		cmd := strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
		switch cmd[0] {
		case "n":
			room <- Cmd{[]string{"n", id}, nil}
			result := <-resp
			switch result.cmd[0] {
			case "locked":
				fmt.Println("The door to the north is locked")
			case "room":
				room = result.resp
				fmt.Println("You enter the room to the north")
				room <- Cmd{[]string{"look", id}, nil}
				result := <-resp
				fmt.Println(result.cmd[0])
			}
		case "e":
			room <- Cmd{[]string{"e", id}, nil}
			result := <-resp
			switch result.cmd[0] {
			case "locked":
				fmt.Println("The door to the east is locked")
			case "room":
				room = result.resp
				fmt.Println("You enter the room to the east")
				room <- Cmd{[]string{"look", id}, nil}
				result := <-resp
				fmt.Println(result.cmd[0])
			}
		case "s":
			room <- Cmd{[]string{"s", id}, nil}
			result := <-resp
			switch result.cmd[0] {
			case "locked":
				fmt.Println("The door to the south is locked")
			case "room":
				room = result.resp
				fmt.Println("You enter the room to the south")
				room <- Cmd{[]string{"look", id}, nil}
				result := <-resp
				fmt.Println(result.cmd[0])
			}
		case "w":
			room <- Cmd{[]string{"w", id}, nil}
			result := <-resp
			switch result.cmd[0] {
			case "locked":
				fmt.Println("The door to the west is locked")
			case "room":
				room = result.resp
				fmt.Println("You enter the room to the west")
				room <- Cmd{[]string{"look", id}, nil}
				result := <-resp
				fmt.Println(result.cmd[0])
			}
		case "look":
			room <- Cmd{[]string{"look", id}, nil}
			result := <-resp
			fmt.Println(result.cmd[0])
		case "take":
			room <- Cmd{[]string{"take", cmd[1], cmd[2]}, nil}
			result := <-resp
			/*
			   case "leave":
			   case "exit":
			*/
		}
	}
}

func main() {
	c0 := make(chan Cmd)
	c1 := make(chan Cmd)
	c2 := make(chan Cmd)
	c3 := make(chan Cmd)
	c4 := make(chan Cmd)
	go room(c1, Doors{nil, c2, c3, nil}, map[string]uint64{"gold": 1}, map[string](chan Cmd){"John Doe": c0})
	go room(c2, Doors{nil, nil, c4, c1}, map[string]uint64{"gold": 2}, map[string](chan Cmd){})
	go room(c3, Doors{c1, c4, nil, nil}, map[string]uint64{"gold": 3}, map[string](chan Cmd){})
	go room(c4, Doors{c2, nil, nil, c3}, map[string]uint64{"gold": 4}, map[string](chan Cmd){})
	player(c0, c1, "John Doe", map[string]uint64{})
}
