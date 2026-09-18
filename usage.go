package main

import (
	"dvc/commands"
	"fmt"
	"sort"
)

func Usage() {
	fmt.Print("Usage: dvc <command> [args]\n\n")

	fmt.Println("Available commands:")

	names := make([]string, 0, len(commands.CommandRegistry))
	for name := range commands.CommandRegistry {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		cmd := commands.CommandRegistry[name]
		fmt.Printf("%s\t%s\n", cmd.Name, cmd.Description)
	}
}
