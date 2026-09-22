package commands

import (
	"dvc/models"
	"flag"
	"fmt"
	"strings"
)

func parseListOptions(args []string) (models.ListRecordsOptions, error) {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	var options models.RecordQueryOptions
	var rawFormat string

	fs.StringVar(&options.Select, "select", "", "Select columns to retrieve and display")
	fs.StringVar(&options.Filter, "filter", "", "Filter rows based on an OData filter condition")
	fs.StringVar(&options.OrderBy, "order-by", "", "Order rows by a specific column")
	fs.StringVar(&options.Expand, "expand", "", "Expand related columns")
	fs.IntVar(&options.Top, "top", 0, "Limit the number of rows returned")
	fs.IntVar(&options.MaxPageSize, "max-page-size", 0, "Limit the number of rows returned per page")
	fs.StringVar(&rawFormat, "output", "table", "Output format (table, json)")
	fs.StringVar(&rawFormat, "o", "table", "Output format (table, json) shorthand")
	fs.Usage = func() { fmt.Println(listUsage()) }

	if err := fs.Parse(args); err != nil {
		return models.ListRecordsOptions{}, err
	}

	if options.Top < 0 {
		return models.ListRecordsOptions{}, fmt.Errorf("top must be a greater than 0")
	}

	if options.MaxPageSize < 0 {
		return models.ListRecordsOptions{}, fmt.Errorf("max-page-size must be a greater than 0")
	}

	var format models.OutputFormat
	switch strings.ToLower(rawFormat) {
	case string(models.OutputFormatJSON):
		format = models.OutputFormatJSON
	case string(models.OutputFormatTable):
		format = models.OutputFormatTable
	default:
		return models.ListRecordsOptions{}, fmt.Errorf("invalid output format %q (allowed: json, table)", rawFormat)
	}

	return models.ListRecordsOptions{Query: options, Output: format}, nil
}
