package commands

import (
	"dvc/client"
	"dvc/models"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var ListCommand = &Command{
	Name:        "list",
	Description: "Retrieve records from a Dataverse table",
	Run:         listRecords,
}

func listRecords(client *client.DataverseClient, args []string) error {
	if len(args) == 0 || slices.Contains(args, "-h") || slices.Contains(args, "--help") {
		fmt.Println(listUsage())
		return nil
	}

	logicalName := args[0]
	optionArgs := args[1:]

	table, err := client.GetTable(logicalName)
	if err != nil {
		return err
	}

	queryOptions, err := parseQueryOptions(optionArgs)
	if err != nil {
		return err
	}

	records, err := client.ListRecords(table.EntitySetName, queryOptions)
	if err != nil {
		return fmt.Errorf("failed to list records: %w", err)
	}

	if len(records) == 0 {
		fmt.Println("no records found")
		return nil
	}

	prettyJSON, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(prettyJSON))
	return nil
}

func parseQueryOptions(args []string) (models.QueryOptions, error) {
	options := models.QueryOptions{}

	parseStringFlag := func(target *string, flag string, i *int) error {
		if *target != "" {
			return fmt.Errorf("duplicate %s", flag)
		}
		if *i+1 >= len(args) {
			return fmt.Errorf("missing value for %s", flag)
		}
		*i++
		*target = args[*i]
		return nil
	}

	parseIntFlag := func(target *int, flag string, i *int) error {
		var valStr string
		if err := parseStringFlag(&valStr, flag, i); err != nil {
			return err
		}
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %w", flag, err)
		}
		if val <= 0 {
			return fmt.Errorf("%s must be greater than 0", flag)
		}
		*target = val
		return nil
	}

	for i := 0; i < len(args); i++ {
		var err error
		switch args[i] {
		case "--select":
			err = parseStringFlag(&options.Select, "--select", &i)
			options.Select = strings.ReplaceAll(options.Select, " ", "")
		case "--filter":
			err = parseStringFlag(&options.Filter, "--filter", &i)
		case "--order-by":
			err = parseStringFlag(&options.OrderBy, "--order-by", &i)
		case "--expand":
			err = parseStringFlag(&options.Expand, "--expand", &i)
		case "--top":
			err = parseIntFlag(&options.Top, "--top", &i)
		case "--max-page-size":
			err = parseIntFlag(&options.MaxPageSize, "--max-page-size", &i)
		default:
			return options, fmt.Errorf("unknown option: %s", args[i])
		}

		if err != nil {
			return options, err
		}
	}

	return options, nil
}

func listUsage() string {
	return `Usage: dvc list <table> [options]

Options:
  -h, --help 			Show this help message
  --select <columns>  		Select specific columns to display (default: all, comma-separated)
  --filter <condition> 		Filter rows based on a condition (default: none)
  --order-by <column> 		Order rows by a specific column
  --expand <columns> 		Expand specific columns to display
  --top <count>     		Limit the number of rows to display
  --max-page-size <count> 	Limit the number of rows per query page (pagination)
`
}
