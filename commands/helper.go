package commands

import (
	"flag"
	"fmt"
	"slices"
)

func checkArgs(args []string, usage string) error {
	return checkArgsWithMinLength(args, 1, usage)
}

func checkArgsWithMinLength(args []string, minLen int, usage string) error {
	containsHelpFlag := func(arg string) bool {
		return arg == "-h" || arg == "--help"
	}

	if slices.ContainsFunc(args, containsHelpFlag) {
		fmt.Println(usage)
		return nil
	} else if len(args) < minLen {
		fmt.Println(usage)
		return flag.ErrHelp
	}

	return nil
}

func valueOrEmpty(record map[string]any, key string) any {
	if value := record[key]; value != nil {
		return value
	}
	return ""
}
