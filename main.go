package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/StelIify/pokedex/internal/data"
	"github.com/chzyer/readline"
)

func cleanInput(text string) []string {
	cleanTxt := strings.TrimSpace(text)
	cleanTxt = strings.ToLower(cleanTxt)
	words := strings.Fields(cleanTxt)
	return words
}

func main() {
	// erorrLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	cache := data.NewCache(time.Minute)
	httpClient := data.NewClient(time.Minute, cache)

	cfg := &Config{
		client:  httpClient,
		pokedex: make(map[string]*data.PokemonResp),
	}
	rl, err := readline.New("pokedex> ")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create readline instance: %v", err))
	}
	defer rl.Close()

	commands := getCommandsMap(getCommands())
	for {
		text, err := rl.Readline()
		if err != nil {
			if errors.Is(readline.ErrInterrupt, err) {
				os.Exit(0)
			} else {
				fmt.Println("Failed to read input: ", err)
				continue
			}
		}
		cleanText := cleanInput(text)
		command, ok := commands[cleanText[0]]
		if !ok {
			fmt.Println("Invalid command, type 'help' to see available commands")
			continue
		}
		args := []string{}
		if len(cleanText) > 1 {
			args = cleanText[1:]
		}
		err = command.callback(cfg, args...)
		if err != nil {
			fmt.Println(err)
			// erorrLog.Println(err)
		}
	}
}
