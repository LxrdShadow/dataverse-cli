package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"dvc/client"
	"dvc/commands"
	"dvc/config"
)

func main() {
	flag.Usage = PrintUsage
	flag.Parse()
	args := os.Args[1:]
	if len(args) == 0 {
		PrintUsage()
		os.Exit(1)
	}

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	// Create a new Dataverse client
	client, err := client.NewDataverseClient(config.BaseURL, config.Token, config.Timeout)
	if err != nil {
		fmt.Println("Error creating client:", err)
		os.Exit(1)
	}

	// Get the selected command
	command := args[0]
	availableCommands := commands.CommandRegistry
	if _, ok := availableCommands[command]; !ok {
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}

	cmd := availableCommands[command]
	err = cmd.Run(client, args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
