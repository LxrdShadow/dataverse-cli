package commands

import (
	"dvc/client"
	"dvc/models"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var ListCommand = &Command{
	Name:        "list",
	Description: "Retrieve records from a Dataverse table",
	Run:         listRecords,
}

func listRecords(client *client.DataverseClient, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: dvc list <table_logical_name> [options]")
	}

	logicalName := args[0]
	optionArgs := args[1:]

	table, err := client.GetTable(logicalName)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	queryOptions, err := parseQueryOptions(optionArgs)
	if err != nil {
		return fmt.Errorf("%w", err)
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

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--select":
			if options.Select != "" {
				return options, fmt.Errorf("duplicate --select")
			}
			if i+1 < len(args) {
				options.Select = strings.ReplaceAll(args[i+1], " ", "")
				i++
			} else {
				return options, fmt.Errorf("missing value for --select")
			}
		case "--filter":
			if options.Filter != "" {
				return options, fmt.Errorf("duplicate --filter")
			}
			if i+1 < len(args) {
				options.Filter = args[i+1]
				i++
			} else {
				return options, fmt.Errorf("missing value for --filters")
			}
		case "--orderby":
			if options.OrderBy != "" {
				return options, fmt.Errorf("duplicate --orderby")
			}
			if i+1 < len(args) {
				options.OrderBy = args[i+1]
				i++
			} else {
				return options, fmt.Errorf("missing value for --orderby")
			}
		case "--expand":
			if options.Expand != "" {
				return options, fmt.Errorf("duplicate --expand")
			}
			if i+1 < len(args) {
				options.Expand = args[i+1]
				i++
			} else {
				return options, fmt.Errorf("missing value for --expand")
			}
		case "--top":
			if options.Top != 0 {
				return options, fmt.Errorf("duplicate --top")
			}
			if i+1 >= len(args) {
				return options, fmt.Errorf("missing value for --top")
			}

			top, err := strconv.ParseInt(args[i+1], 10, 0)
			if err != nil {
				return options, fmt.Errorf("invalid value for --top: %w", err)
			}

			if top <= 0 {
				return options, fmt.Errorf("--top must be greater than 0")
			}

			options.Top = int(top)
			i++
		case "--max-page-size":
			if options.MaxPageSize != 0 {
				return options, fmt.Errorf("duplicate --max-page-size")
			}
			if i+1 >= len(args) {
				return options, fmt.Errorf("missing value for --max-page-size")
			}

			maxPageSize, err := strconv.ParseInt(args[i+1], 10, 0)
			if err != nil {
				return options, fmt.Errorf("invalid value for --max-page-size: %w", err)
			}

			if maxPageSize <= 0 {
				return options, fmt.Errorf("--max-page-size must be greater than 0")
			}

			options.MaxPageSize = int(maxPageSize)
			i++

		default:
			return options, fmt.Errorf("unknown option: %s", args[i])
		}
	}

	return options, nil
}
