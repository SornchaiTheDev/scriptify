package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Command struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

type Config struct {
	Commands []Command `json:"commands"`
}

const configFileName = ".scriptify.json"

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

func loadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	config := &Config{Commands: []Command{}}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func saveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func addCommand(name, command string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	for i, cmd := range config.Commands {
		if cmd.Name == name {
			config.Commands[i].Command = command
			return saveConfig(config)
		}
	}

	config.Commands = append(config.Commands, Command{Name: name, Command: command})
	return saveConfig(config)
}

func executeCommand(name string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	for _, cmd := range config.Commands {
		if cmd.Name == name {
			parts := strings.Fields(cmd.Command)
			if len(parts) == 0 {
				return fmt.Errorf("empty command for '%s'", name)
			}

			command := exec.Command(parts[0], parts[1:]...)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			command.Stdin = os.Stdin
			return command.Run()
		}
	}

	return fmt.Errorf("command '%s' not found", name)
}

func showHelp() {
	fmt.Println("scriptify - A simple command proxy tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  scriptify add <name> <command>    Add or update a command")
	fmt.Println("  scriptify help                    Show this help message")
	fmt.Println("  scriptify <name>                  Execute a stored command")
	fmt.Println()

	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	if len(config.Commands) > 0 {
		fmt.Println("Available commands:")
		for _, cmd := range config.Commands {
			fmt.Printf("  %-20s %s\n", cmd.Name, cmd.Command)
		}
	} else {
		fmt.Println("No commands configured. Use 'scriptify add <name> <command>' to add one.")
	}
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: scriptify add <name> <command>")
			os.Exit(1)
		}
		name := os.Args[2]
		command := strings.Join(os.Args[3:], " ")

		err := addCommand(name, command)
		if err != nil {
			fmt.Printf("Error adding command: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Command '%s' added successfully\n", name)

	case "help":
		showHelp()

	default:
		commandName := os.Args[1]
		err := executeCommand(commandName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}