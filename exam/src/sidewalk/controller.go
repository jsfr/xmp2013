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
func Controller(clockCtrl chan ClockRequest, injCtrl [](chan int),
	drunkCtrl chan Coordinate, printer chan Coordinate, cells [][](chan Msg),
	x int, y int) {
	os.Stdout.WriteString("\033[0;0H\033[J\033[" + strconv.Itoa(y+6) + ";0H")
	clockCtrl <- ClockRequest{CmdResume, nil}
	play := true
	exit := false
	os.Stdout.WriteString("\r\n" + HelpText + "\033[F")
	for !exit {
		os.Stdout.WriteString("\033[Kcmd> ")
		stdin := bufio.NewReader(os.Stdin)
		input, _ := stdin.ReadString('\n')
		cmd := strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
		os.Stdout.WriteString("\033[J")
		switch cmd[0] {
		case "pause":
			if play {
				clockCtrl <- ClockRequest{CmdPause, nil}
				os.Stdout.WriteString("Simulation paused")
				play = !play
			} else {
				os.Stdout.WriteString("Simulation already paused")
			}
		case "resume":
			if !play {
				clockCtrl <- ClockRequest{CmdResume, nil}
				os.Stdout.WriteString("Simulation resumed")
				play = !play
			} else {
				os.Stdout.WriteString("Simulation already running")
			}
		case "rate":
			i, _ := strconv.Atoi(cmd[1])
			if i > 0 && i < 11 {
				for _, inj := range injCtrl {
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
				drunkCtrl <- Coordinate{drunkX, drunkY}
				msg := <-drunkCtrl
				if msg.X == drunkX && msg.Y == drunkY {
					os.Stdout.WriteString("Drunk placed at (" + cmd[1] + "," +
						cmd[2] + ")")
				} else {
					os.Stdout.WriteString("Drunk not placed, coordinates are" +
						" occupied or drunk already exists on board")
				}
			} else {
				os.Stdout.WriteString("Drunk not placed, coordinates are invalid")
			}
		case "timestep":
			wait, _ := strconv.Atoi(cmd[1])
			if wait >= 0 {
				clockCtrl <- ClockRequest{CmdTime, []int{wait}}
				os.Stdout.WriteString("Time between steps changed")
			} else {
				os.Stdout.WriteString("Time between steps cannot be negative")
			}
		case "help":
			os.Stdout.WriteString(HelpText)
		case "exit":
			clockCtrl <- ClockRequest{CmdExit, nil}
			for _, i := range cells {
				for _, j := range i {
					close(j)
				}
			}
			for _, i := range injCtrl {
				close(i)
			}
			close(drunkCtrl)
			close(printer)
			close(clockCtrl)
			return
		default:
			os.Stdout.WriteString("Unknown command")
		}
		os.Stdout.WriteString("\033[F")
	}
}

const HelpText = "The following commands can be called\r\n" +
	"pause\033[25Gpauses the simulation after the current step has executed\r\n" +
	"resume\033[25Gresumes the simulation if it has been paused\r\n" +
	"rate [ARG]\033[25Gsets the number of steps between injections of pedestrians\r\n" +
	"drunk [ARG_X] [ARG_Y]\033[25Ginsert a drunk at the given position\r\n" +
	"timestep [ARG]\033[25Gchanges the time (in milliseconds) " +
	"simulation waits between step\r\n" +
	"help\033[25Gshows this text\r\n" +
	"exit\033[25Gexits the simulation\033[7F"
