package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
	config      *Config
}

type ApiResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func generateLocationConfig() *Config {
	return &Config{
		next:     "https://pokeapi.co/api/v2/location-area/",
		previous: "",
	}
}

func generatePokemonConfig() *Config {
	return &Config{
		next:     "https://pokeapi.co/api/v2/pokemon/",
		previous: "",
	}
}

func getCommands() []cliCommand {
	mapConfig := generateLocationConfig()
	pokemonConfig := generatePokemonConfig()

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
			config:      mapConfig,
		},
		{
			name:        "mapb",
			description: "Dispays the name of 20 previous location areas",
			callback:    commandMapBack,
			config:      mapConfig,
		},
		{
			name:        "pokemon",
			description: "Dispays the name of 20 pokemons",
			callback:    commandPokemon,
			config:      pokemonConfig,
		},
		{
			name:        "pokemonb",
			description: "Dispays the name of 20 prevoius pokemons",
			callback:    commandPokemonBack,
			config:      pokemonConfig,
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

func commandHelp(c *Config) error {
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

func commandExit(c *Config) error {
	os.Exit(0)
	return nil
}

func fetchData(url string) (*ApiResponse, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error during get request to the next batch of data: %v", err)
	}
	var apiResponse ApiResponse
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("error during response unmarshal: %v", err)
	}

	return &apiResponse, nil
}

func commandMap(c *Config) error {
	apiResponse, err := fetchData(c.next)

	if err != nil {
		return fmt.Errorf("commandMap: %v", err)
	}
	c.next = apiResponse.Next
	c.previous = apiResponse.Previous

	fmt.Println("20 next location area names:")
	for _, item := range apiResponse.Results {
		fmt.Println(item.Name)
	}

	return nil
}

func commandMapBack(c *Config) error {
	if c.previous == "" {
		return fmt.Errorf("you are on the first page, there is no previous locations")
	}
	apiResponse, err := fetchData(c.previous)
	if err != nil {
		return fmt.Errorf("commandMapBack: %v", err)
	}
	c.previous = apiResponse.Previous
	c.next = apiResponse.Next

	fmt.Println("20 previous location area names:")
	for _, item := range apiResponse.Results {
		fmt.Println(item.Name)
	}

	return nil
}

func commandPokemon(c *Config) error {
	apiResponse, err := fetchData(c.next)

	if err != nil {
		return fmt.Errorf("commandPokemon: %v", err)
	}
	c.next = apiResponse.Next
	c.previous = apiResponse.Previous

	fmt.Println("20 next pokemon names:")
	for _, item := range apiResponse.Results {
		fmt.Println(item.Name)
	}

	return nil
}

func commandPokemonBack(c *Config) error {
	if c.previous == "" {
		return fmt.Errorf("you are on the first page, there is no previous pokemons")
	}
	apiResponse, err := fetchData(c.previous)

	if err != nil {
		return fmt.Errorf("commandPokemonBack: %v", err)
	}
	c.previous = apiResponse.Previous
	c.next = apiResponse.Next

	fmt.Println("20 previous pokemon names:")
	for _, item := range apiResponse.Results {
		fmt.Println(item.Name)
	}

	return nil
}
