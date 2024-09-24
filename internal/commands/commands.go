package commands

import (
	"fmt"
	"github.com/znataniel/go-pokedex/internal/pokeapi"
	"os"
)

type CommandFn func(*Config) error

type Command struct {
	name     string
	desc     string
	Callback CommandFn
}

type Config struct {
	Commands    map[string]Command
	UserCommand string
	MapConfig   pokeapi.PokeMapConfig
}

func InitializeCommands() map[string]Command {
	commands := make(map[string]Command)
	commands["help"] = Command{
		name:     "help",
		desc:     "Shows available commands",
		Callback: commandHelp,
	}

	commands["exit"] = Command{
		name:     "exit",
		desc:     "Exits the session",
		Callback: commandExit,
	}

	commands["map"] = Command{
		name:     "map",
		desc:     "Traverses the map forwards",
		Callback: commandMap,
	}

	commands["mapb"] = Command{
		name:     "mapb",
		desc:     "Traverses the map backwards",
		Callback: commandMapb,
	}

	return commands
}

func commandHelp(config *Config) error {
	commandMap := config.Commands
	for _, comm := range commandMap {
		fmt.Println(comm.name, "---", comm.desc)
	}
	return nil
}

func commandExit(config *Config) error {
	os.Exit(0)
	return nil
}

func commandMap(config *Config) error {
	return pokeapi.PkmnMap(&(*config).MapConfig)
}

func commandMapb(config *Config) error {
	return pokeapi.PkmnMapb(&(*config).MapConfig)
}
