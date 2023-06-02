package main

import "testing"

func TestCliCommands(t *testing.T) {
	commands := getCommands()

	expectedCommands := []cliCommands{
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
	}

	expectedLength := len(expectedCommands)

	for i := 0; i < expectedLength; i++ {
		if commands[i].name != expectedCommands[i].name {
			t.Errorf("Expected %s got %s", expectedCommands[i].name, commands[i].name)
		}
		if commands[i].description != expectedCommands[i].description {
			t.Errorf("Expected %s got %s", expectedCommands[i].description, commands[i].description)
		}
		if commands[i].callback == nil {
			t.Errorf("Expected non-nil callback for the command")
		}
	}

}

func TestGetCommandsMap(t *testing.T) {
	commands := getCommands()
	commandMap := getCommandsMap(commands)
	for _, command := range commands {
		_, ok := commandMap[command.name]
		if !ok {
			t.Errorf("Command with name '%s' is not found in the command map", command.name)
			continue
		}
	}
}

func TestCleanInput(t *testing.T) {
	testString := "Help "

	got := cleanInput(testString)
	expected := "help"

	if got != expected {
		t.Errorf("Got %s, expected %s", got, expected)
	}
}
