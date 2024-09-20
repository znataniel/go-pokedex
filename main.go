package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	commands := initializeCommands()
	for {
		fmt.Print("Input > ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan() //scans ONCE till EOL or error
		if scanner.Err() != nil {
			fmt.Println("Error reading from stdin")
			continue
		}

		line := scanner.Text()
		if line == "" {
			continue
		}
		_, ok := commands[line]
		if !ok {
			fmt.Println("Command not found:", line)
			continue
		}

		if line == "exit" {
			break
		}

		commands[line].callback(commands)
	}
}
