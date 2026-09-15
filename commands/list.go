package commands

import (
	"dvc/client"
	"dvc/models"
	"encoding/json"
	"fmt"
	"strconv"
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
			if i+1 < len(args) {
				options.Select = args[i+1]
				i++
			} else {
				return options, fmt.Errorf("missing value for --select")
			}
		case "--filter":
			if i+1 < len(args) {
				options.Filter = args[i+1]
				i++
			} else {
				return options, fmt.Errorf("missing value for --filters")
			}
		case "--orderby":
			if i+1 < len(args) {
				options.OrderBy = args[i+1]
				i++
			} else {
				return options, fmt.Errorf("missing value for --orderby")
			}
		case "--expand":
			if i+1 < len(args) {
				options.Expand = args[i+1]
				i++
			} else {
				return options, fmt.Errorf("missing value for --expand")
			}
		case "--top":
			if i+1 < len(args) {
				top, err := strconv.ParseInt(args[i+1], 10, 0)
				if err != nil {
					return options, fmt.Errorf("invalid value for --top: %w", err)
				}
				options.Top = int(top)
				i++
			} else {
				return options, fmt.Errorf("missing value for --top")
			}
		default:
			return options, fmt.Errorf("unknown option: %s", args[i])
		}
	}

	return options, nil
}
