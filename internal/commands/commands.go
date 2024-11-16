package commands

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/znataniel/go-pokedex/internal/pokeapi"
)

type CommandFn func(*Config) error

type Command struct {
	name     string
	desc     string
	Callback CommandFn
}

type Config struct {
	Commands         map[string]Command
	UserCommand      string
	CommandArguments string
	MapConfig        pokeapi.PokeMapConfig
	Pokedex          map[string]pokeapi.Pokemon
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

	commands["explore"] = Command{
		name:     "explore <area-name>",
		desc:     "Explores a given area",
		Callback: commandExplore,
	}

	commands["catch"] = Command{
		name:     "catch <pokemon-name>",
		desc:     "Tries to catch a pokemon!",
		Callback: commandCatch,
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
	cfg := &config.MapConfig
	url := pokeapi.BaseUrl + "location-area/?offset=0&limit=20" // First call for locations

	if cfg.Next != nil {
		url = *cfg.Next
	}

	pokeRes, err := pokeapi.GetLocations(url)
	if err != nil {
		return err
	}

	for _, loc := range pokeRes.Results {
		fmt.Println(loc.Name)
	}

	cfg.Prev = pokeRes.Previous
	cfg.Next = pokeRes.Next

	return nil
}

func commandMapb(config *Config) error {
	cfg := &config.MapConfig

	if cfg.Prev == nil {
		return fmt.Errorf("error: no previous locations available")
	}

	url := *cfg.Prev

	pokeRes, err := pokeapi.GetLocations(url)
	if err != nil {
		return err
	}

	for _, loc := range pokeRes.Results {
		fmt.Println(loc.Name)
	}

	cfg.Prev = pokeRes.Previous
	cfg.Next = pokeRes.Next

	return nil
}

func commandExplore(config *Config) error {
	if config.CommandArguments == "" {
		return fmt.Errorf("error: no area to explore provided")
	}

	url := pokeapi.BaseUrl + "location-area/" + config.CommandArguments

	pokeRes, err := pokeapi.GetLocationData(url)
	if err != nil {
		return err
	}

	for _, encounter := range pokeRes.PokemonEncounters {
		fmt.Println("\t- ", encounter.Pokemon.Name)
	}

	return nil

}

func commandCatch(config *Config) error {
	if config.CommandArguments == "" {
		return fmt.Errorf("error: no pokemon to catch provided")
	}

	url := pokeapi.BaseUrl + "pokemon/" + config.CommandArguments
	pokeRes, err := pokeapi.GetPokemon(url)
	if err != nil {
		return err
	}

	fmt.Println("throwing pokeball at", config.CommandArguments)

	probability := math.Floor(100 - 3.85*math.Sqrt(float64(pokeRes.BaseExperience)))
	pick := rand.Intn(100)

	for i := 0; i < 3; i++ {
		time.Sleep((2000 / 3) * time.Millisecond)
		fmt.Print(".")
	}
	time.Sleep((2000 / 3) * time.Millisecond)
	fmt.Println()

	if pick > int(probability) {
		fmt.Println("oh no", config.CommandArguments, "escaped!")
		return nil
	}

	fmt.Println("yes!", config.CommandArguments, "has been catched")
	_, exists := config.Pokedex[config.CommandArguments]
	if !exists {
		config.Pokedex[config.CommandArguments] = pokeRes
	}
	return nil

}
