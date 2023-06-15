package main

import (
	"errors"
	"fmt"
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
