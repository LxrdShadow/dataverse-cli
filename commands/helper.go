package commands

import (
	"dvc/errors"
	"flag"
	"fmt"
	"slices"
)

func validateArgs(args []string, minimum int, usage string) error {
	if slices.Contains(args, "-h") || slices.Contains(args, "--help") {
		fmt.Println(usage)
		return flag.ErrHelp
	}

	if len(args) < minimum {
		fmt.Println(usage)
		return errors.ErrUsage
	}

	return nil
}

func valueOrEmpty(record map[string]any, key string) any {
	if value := record[key]; value != nil {
		return value
	}
	return ""
}
