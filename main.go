package main

import (
	"bufio"
	"fmt"
	"github.com/znataniel/go-pokedex/internal/commands"
	"github.com/znataniel/go-pokedex/internal/pokeapi"
	"os"
	"strings"
)

func main() {
	commsMap := commands.InitializeCommands()
	arguments := commands.Config{
		Commands:    commsMap,
		UserCommand: "",
		MapConfig:   pokeapi.PokeMapConfig{},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")

		scanner.Scan() //scans ONCE till EOL or error
		if scanner.Err() != nil {
			fmt.Println("Error reading from stdin")
			continue
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		command, arg, _ := strings.Cut(line, " ")

		_, ok := commsMap[command]
		if !ok {
			fmt.Println("Command not found:", command)
			continue
		}

		arguments.UserCommand = command
		arguments.CommandArguments = arg

		if err := commsMap[command].Callback(&arguments); err != nil {
			fmt.Println(err)
		}
	}
}
