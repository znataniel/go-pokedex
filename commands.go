package main

import "fmt"

type CommandFn func(map[string]Command)

type Command struct {
	name     string
	desc     string
	callback CommandFn
}

func commandHelp(commands map[string]Command) {
	for _, comm := range commands {
		fmt.Println(comm.name, "---", comm.desc)
	}
}

func initializeCommands() map[string]Command {
	commands := make(map[string]Command)
	commands["help"] = Command{
		name:     "help",
		desc:     "Shows available commands",
		callback: commandHelp,
	}

	commands["exit"] = Command{
		name:     "exit",
		desc:     "exits the session",
		callback: func(map[string]Command) {},
	}

	return commands
}
