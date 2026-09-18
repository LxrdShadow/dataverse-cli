package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"dvc/client"
	"dvc/commands"
	"dvc/config"
)

func main() {
	flag.Usage = Usage
	flag.Parse()
	args := os.Args[1:]
	if len(args) == 0 {
		Usage()
		return
	}

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	config, err := config.Load()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Create a new Dataverse client
	client := client.NewDataverseClient(config.BaseURL, config.Token, config.Timeout)

	// Get the selected command
	command := args[0]
	availableCommands := commands.CommandRegistry
	if _, ok := availableCommands[command]; !ok {
		fmt.Println("Unknown command:", command)
		return
	}

	cmd := availableCommands[command]
	err = cmd.Run(client, args[1:])
	if err != nil {
		fmt.Println("Error:", err)
	}
}
