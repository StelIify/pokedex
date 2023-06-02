package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommands struct {
	name        string
	description string
	callback    func(*Config) error
}

func getCommands() []cliCommands {
	commands := []cliCommands{
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
	}
	return commands
}
func getCommandsMap(commands []cliCommands) map[string]cliCommands {
	commandMap := make(map[string]cliCommands)
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

	return nil
}

func commandExit(c *Config) error {
	os.Exit(0)
	return nil
}

type ApiResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func commandMapBack(c *Config) error {
	if c.previous == "" {
		return fmt.Errorf("first page, there is no previous locations")
	}
	response, err := http.Get(c.previous)
	if err != nil {
		return err
	}
	var apiResponse ApiResponse
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	_ = json.Unmarshal(body, &apiResponse)
	c.previous = apiResponse.Previous
	c.next = apiResponse.Next

	for _, item := range apiResponse.Results {
		fmt.Println(item.Name)
	}

	return nil
}
func commandMap(c *Config) error {
	response, err := http.Get(c.next)
	if err != nil {
		return fmt.Errorf("error during get request to the next")
	}
	var apiResponse ApiResponse
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	_ = json.Unmarshal(body, &apiResponse)

	c.next = apiResponse.Next
	c.previous = apiResponse.Previous

	for _, item := range apiResponse.Results {
		fmt.Println(item.Name)
	}

	return nil
}
