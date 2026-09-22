package commands

import (
	"flag"
	"fmt"
	"strings"

	"dvc/models"
)

func parseGetOptions(args []string) (models.GetOptions, error) {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)
	var options models.GetOptions
	var rawFormat string

	fs.StringVar(&options.Select, "select", "", "Select columns to retrieve and display")
	fs.StringVar(&rawFormat, "output", "table", "Output format (table, json)")
	fs.StringVar(&rawFormat, "o", "table", "Output format (table, json) shorthand")

	if err := fs.Parse(args); err != nil {
		return models.GetOptions{}, err
	}

	var format models.OutputFormat
	switch strings.ToLower(rawFormat) {
	case string(models.OutputFormatJSON):
		format = models.OutputFormatJSON
	case string(models.OutputFormatTable):
		format = models.OutputFormatTable
	default:
		return models.GetOptions{}, fmt.Errorf("invalid output format %q (allowed: json, table)", rawFormat)
	}

	return models.GetOptions{
		Select: options.Select,
		Output: format,
	}, nil
}
