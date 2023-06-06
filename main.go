package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	next     string
	previous string
}

// type CommandConfig struct {
// 	configs map[string]*Config
// }

// func NewCommandConfig() *CommandConfig {
// 	return &CommandConfig{
// 		configs: map[string]*Config{
// 			"map": generateLocationConfig(),
// 			"pokemon": generatePokemonConfig(),
// 		},

// 	}
// }

func NewConfig(next, previous string) *Config {
	return &Config{
		next:     next,
		previous: previous,
	}
}

const cliName = "pokedex"

func printPrompt() {
	fmt.Print(cliName, "> ")
}

func cleanInput(text string) string {
	cleanTxt := strings.TrimSpace(text)
	cleanTxt = strings.ToLower(cleanTxt)
	return cleanTxt
}
func main() {
	erorrLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	reader := bufio.NewScanner(os.Stdin)
	printPrompt()
	commands := getCommandsMap(getCommands())
	for reader.Scan() {
		text := reader.Text()
		cleanText := cleanInput(text)
		command, ok := commands[cleanText]
		if !ok {
			fmt.Println("Invalid command, type 'help' to see available commands")
			continue
		}
		err := command.callback(command.config)
		if err != nil {
			fmt.Println(err)
			erorrLog.Println(err)
		}
		printPrompt()
	}
}
