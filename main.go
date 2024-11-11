package main

import (
	"bufio"
	"fmt"
	"github.com/znataniel/go-pokedex/internal/commands"
	"github.com/znataniel/go-pokedex/internal/pokeapi"
	"os"
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

		line := scanner.Text()
		if line == "" {
			continue
		}

		_, ok := commsMap[line]
		if !ok {
			fmt.Println("Command not found:", line)
			continue
		}

		arguments.UserCommand = line

		if err := commsMap[line].Callback(&arguments); err != nil {
			fmt.Println(err)
		}
	}
}
