package sidewalk

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func Controller(clock chan bool, injs [](chan int), y int) {
	os.Stdout.WriteString("\033[0;0H\033[J\033[" + strconv.Itoa(y+4) + ";0H")
	clock <- true
	for {
		os.Stdout.WriteString("\033[Kcmd> ")
		stdin := bufio.NewReader(os.Stdin)
		input, _ := stdin.ReadString('\n')
		cmd := strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
		switch cmd[0] {
		case "pause":
			clock <- false
			os.Stdout.WriteString("Simulation paused\033[F")
		case "resume":
			clock <- true
			os.Stdout.WriteString("Simulation resumed\033[F")
		case "rate":
			i, _ := strconv.Atoi(cmd[1])
			if i > 0 && i < 11 {
				for _, inj := range injs {
					inj <- i
				}
				os.Stdout.WriteString("Rate changed to " + cmd[1] + "\033[F")
			} else {
				os.Stdout.WriteString("Rate but be between 1 and 10\033[F")
			}
		case "drunk":
		case "exit":
		default:
			os.Stdout.WriteString("Unknown command\033[F")
		}
	}
}
