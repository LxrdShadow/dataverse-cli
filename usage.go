package main

import (
	"dvc/commands"
	"fmt"
)

func Usage() {
	fmt.Print("Usage: dvc <command> [args]\n\n")

	fmt.Println("Available commands:")
	for _, cmd := range commands.CommandRegistry {
		fmt.Printf("%s\t%s\n", cmd.Name, cmd.Description)
	}
}
