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

	commands["inspect"] = Command{
		name:     "inspect <pokemon-name>",
		desc:     "Inspects a pokemon from the pokedex",
		Callback: commandInspect,
	}

	commands["pokedex"] = Command{
		name:     "pokedex",
		desc:     "Prints all caught pokemon",
		Callback: commandPokedex,
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

	fmt.Println("yes!", config.CommandArguments, "has been caught")
	_, exists := config.Pokedex[config.CommandArguments]
	if !exists {
		config.Pokedex[config.CommandArguments] = pokeRes
	}
	return nil

}

func commandInspect(config *Config) error {
	if config.CommandArguments == "" {
		return fmt.Errorf("error: no pokemon to catch provided")
	}

	pokemon, exists := config.Pokedex[config.CommandArguments]
	if !exists {
		fmt.Println("You haven't caught that pokemon!")
		return nil
	}

	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")
	for _, el := range pokemon.Stats {
		fmt.Println("  -", el.Stat.Name+":", el.BaseStat)
	}
	fmt.Println("Types:")
	for _, el := range pokemon.Types {
		fmt.Println("  -", el.Type.Name)
	}

	return nil
}

func commandPokedex(config *Config) error {
	if n := len(config.Pokedex); n == 0 {
		fmt.Println("your pokedex is empty!")
		return nil
	}

	for pokemon := range config.Pokedex {
		fmt.Println("  -", pokemon)
	}
	return nil
}
