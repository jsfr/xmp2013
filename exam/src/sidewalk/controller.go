package sidewalk

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

/**
 * The shell used for interaction with the player. The shell expects ANSI
 * escape characters to work in the terminal, as these are used to overwrite
 * old output.
 */
func Controller(clock chan ClockRequest, injs [](chan int), drunk chan Coordinate, x int, y int) {
	os.Stdout.WriteString("\033[0;0H\033[J\033[" + strconv.Itoa(y+6) + ";0H")
	clock <- ClockRequest{CmdResume, nil}
	play := true
	for {
		os.Stdout.WriteString("\033[Kcmd> ")
		stdin := bufio.NewReader(os.Stdin)
		input, _ := stdin.ReadString('\n')
		cmd := strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
		os.Stdout.WriteString("\033[K")
		switch cmd[0] {
		case "pause":
			if play {
				clock <- ClockRequest{CmdPause, nil}
				os.Stdout.WriteString("Simulation paused")
				play = !play
			} else {
				os.Stdout.WriteString("Simulation already paused")
			}
		case "resume":
			if !play {
				clock <- ClockRequest{CmdResume, nil}
				os.Stdout.WriteString("Simulation resumed")
				play = !play
			} else {
				os.Stdout.WriteString("Simulation already running")
			}
		case "rate":
			i, _ := strconv.Atoi(cmd[1])
			if i > 0 && i < 11 {
				for _, inj := range injs {
					inj <- i
				}
				os.Stdout.WriteString("Rate changed to " + cmd[1])
			} else {
				os.Stdout.WriteString("Rate but be between 1 and 10")
			}
		case "drunk":
			drunkX, _ := strconv.Atoi(cmd[1])
			drunkY, _ := strconv.Atoi(cmd[2])
			if drunkX >= 0 && drunkX < x && drunkY >= 0 && drunkY < y {
				drunk <- Coordinate{drunkX, drunkY}
				msg := <-drunk
				if msg.X == drunkX && msg.Y == drunkY {
					os.Stdout.WriteString("Drunk placed")
				} else {
					os.Stdout.WriteString("Drunk not placed, coordinates are occupied or drunk already exists on board")
				}
			} else {
				os.Stdout.WriteString("Drunk not placed, coordinates are invalid")
			}
		case "timestep":
			wait, _ := strconv.Atoi(cmd[1])
			if wait >= 0 {
				clock <- ClockRequest{CmdTime, []int{wait}}
				os.Stdout.WriteString("Time between steps changed")
			} else {
				os.Stdout.WriteString("Time between steps cannot be negative")
			}
		case "exit":
		default:
			os.Stdout.WriteString("Unknown command")
		}
		os.Stdout.WriteString("\033[F")
	}
}
