package commands

import (
	"flag"
	"fmt"
	"strings"

	"dvc/models"
)

func parseGetRecordOptions(args []string) (models.GetRecordOptions, error) {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)
	var rawFormat, selectFields string
	var includeEmpty bool

	fs.StringVar(&selectFields, "select", "", "Select columns to retrieve and display")
	fs.BoolVar(&includeEmpty, "include-empty", false, "Include empty fields")
	fs.StringVar(&rawFormat, "output", "table", "Output format (table, json)")
	fs.StringVar(&rawFormat, "o", "table", "Output format (table, json) shorthand")
	fs.Usage = func() { fmt.Println(getUsage()) }

	if err := fs.Parse(args); err != nil {
		return models.GetRecordOptions{}, err
	}

	var format models.OutputFormat
	switch strings.ToLower(rawFormat) {
	case string(models.OutputFormatJSON):
		format = models.OutputFormatJSON
	case string(models.OutputFormatTable):
		format = models.OutputFormatTable
	default:
		return models.GetRecordOptions{}, fmt.Errorf("invalid output format %q (allowed: json, table)", rawFormat)
	}

	return models.GetRecordOptions{
		Select:       selectFields,
		Output:       format,
		IncludeEmpty: includeEmpty,
	}, nil
}
