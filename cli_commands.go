package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/StelIify/pokedex/internal/data"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, ...string) error
}

type Config struct {
	client   data.Client
	pokedex  map[string]*data.PokemonResp
	next     *string
	previous *string
}

func getCommands() []cliCommand {
	commands := []cliCommand{
		{
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		{
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		{
			name:        "map",
			description: "Dispays the name of 20 location areas",
			callback:    commandMap,
		},
		{
			name:        "mapb",
			description: "Dispays the name of 20 previous location areas",
			callback:    commandMapBack,
		},
		{
			name:        "explore",
			description: "Explore particular area and display found pokemons",
			callback:    commandExplore,
		},
		{
			name:        "catch",
			description: "Try to catch pokemon found in particular location",
			callback:    commandCatch,
		},
		{
			name:        "inspect",
			description: "Inspect pokemon in your pokedex, show it's stats",
			callback:    commandInspect,
		},
		{
			name:        "pokedex",
			description: "Display all the pokemons you were able to catch",
			callback:    commandPokedex,
		},
	}
	return commands
}

func getCommandsMap(commands []cliCommand) map[string]cliCommand {
	commandMap := make(map[string]cliCommand)
	for _, command := range commands {
		commandMap[command.name] = command
	}
	return commandMap
}

func commandHelp(c *Config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	fmt.Println()

	commands := getCommands()
	for _, command := range commands {
		fmt.Printf(" - %s: %s\n", command.name, command.description)
	}
	fmt.Println()

	return nil
}

func commandExit(c *Config, args ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(c *Config, args ...string) error {
	apiResponse, err := c.client.GetLocationAreas(c.next)

	if err != nil {
		return fmt.Errorf("commandMap: %v", err)
	}
	fmt.Println("20 next location area names:")
	for _, item := range apiResponse.Results {
		fmt.Printf(" - %s\n", item.Name)
	}
	c.next = apiResponse.Next
	c.previous = apiResponse.Previous

	return nil
}

func commandMapBack(c *Config, args ...string) error {
	if c.previous == nil {
		return fmt.Errorf("you are on the first page, there is no previous locations")
	}
	apiResponse, err := c.client.GetLocationAreas(c.previous)
	if err != nil {
		return fmt.Errorf("commandMapBack: %v", err)
	}
	fmt.Println("20 previous location area names:")
	for _, item := range apiResponse.Results {
		fmt.Printf(" - %s\n", item.Name)
	}
	c.previous = apiResponse.Previous
	c.next = apiResponse.Next

	return nil
}

func commandExplore(c *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("invalid input format. Usage: explore <location>")
	}
	locationArea := args[0]
	response, err := c.client.GetLocationAreaDetails(locationArea)
	if err != nil {
		return fmt.Errorf("commandExplore: %v", err)
	}
	fmt.Printf("Exploring %s...\n", locationArea)
	fmt.Println("Found Pokemons: ")
	for _, item := range response.PokemonEncounters {
		fmt.Printf("- %s\n", item.Pokemon.Name)
	}
	return nil
}
func generateCatchChance(baseExperience int) float64 {
	maxBaseExp := 400 // max exp value for the pokemon

	catchChance := 1 - (float64(baseExperience) / float64(maxBaseExp))

	if catchChance < 0 {
		catchChance = 0
	} else if catchChance > 1 {
		catchChance = 1
	}

	return catchChance
}

func commandCatch(c *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("invalid input format. Usage: catch <pokemon name>")
	}
	pokemonName := args[0]
	response, err := c.client.GetPokemonDetails(pokemonName)
	if err != nil {
		return fmt.Errorf("commandCatch: %v", err)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	catchChance := generateCatchChance(response.BaseExperience)
	if catchChance > 0 && rand.Float64() > catchChance {
		return fmt.Errorf("%s has escaped, you can try to catch it again", pokemonName)
	}
	_, ok := c.pokedex[pokemonName]
	if ok {
		return fmt.Errorf("you already have %s in your pokedex, try to inspect it", pokemonName)
	}
	fmt.Printf("You have caught %s! You can inspect it now\n", pokemonName)
	c.pokedex[pokemonName] = response
	return nil
}

func commandInspect(c *Config, args ...string) error {
	if len(args) != 1 {
		return errors.New("invalid input format. Usage: catch <pokemon name>")
	}
	pokemonName := args[0]
	pokemon, ok := c.pokedex[pokemonName]
	if !ok {
		return errors.New("you have not caught that pokemon yet")
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Species: %s\n", pokemon.Species.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Base Experience: %d\n", pokemon.BaseExperience)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, ptype := range pokemon.Types {
		fmt.Printf("  - %s\n", ptype.Type.Name)
	}

	return nil
}

func commandPokedex(c *Config, args ...string) error {
	if len(c.pokedex) == 0 {
		return errors.New("you have not caught any pokemon yet, you have much to explore")
	}
	fmt.Println("Your Pokedex:")
	for key, _ := range c.pokedex {
		fmt.Printf(" - %s\n", key)
	}
	return nil
}
