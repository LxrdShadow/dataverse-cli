package commands

import (
	"dvc/models"
	"flag"
	"fmt"
	"strings"
)

func parseEntityMetadataOptions(args []string) (models.EntityMetadataOptions, error) {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)
	var attributes bool
	var rawFormat string

	fs.BoolVar(&attributes, "columns", false, "Display the metadata of the attributes of the entity")
	fs.StringVar(&rawFormat, "output", "table", "Output format (table, json)")
	fs.StringVar(&rawFormat, "o", "table", "Output format (table, json) shorthand")
	fs.Usage = func() { fmt.Println(metadataUsage()) }

	if err := fs.Parse(args); err != nil {
		return models.EntityMetadataOptions{}, err
	}

	var format models.OutputFormat
	switch strings.ToLower(rawFormat) {
	case string(models.OutputFormatJSON):
		format = models.OutputFormatJSON
	case string(models.OutputFormatTable):
		format = models.OutputFormatTable
	default:
		return models.EntityMetadataOptions{}, fmt.Errorf("invalid output format %q (allowed: json, table)", rawFormat)
	}

	return models.EntityMetadataOptions{
		Output:     format,
		Attributes: attributes,
	}, nil
}
