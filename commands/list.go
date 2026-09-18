package commands

import (
	"dvc/client"
	"dvc/models"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
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

	entity, err := client.GetTable(logicalName, true)
	if err != nil {
		return err
	}

	listOptions, err := parseOptions(optionArgs)
	if err != nil {
		return err
	}

	records, err := client.ListRecords(entity.EntitySetName, listOptions)
	if err != nil {
		return fmt.Errorf("failed to list records: %w", err)
	}

	if len(records) == 0 {
		fmt.Println("no records found")
		return nil
	}

	switch listOptions.Output {
	case models.OutputFormatJSON:
		prettyJSON, err := json.MarshalIndent(records, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(prettyJSON))
	case models.OutputFormatTable:
		attributes := entity.Attributes
		// prettyJSON, err := json.MarshalIndent(attributes, "", "  ")
		// if err != nil {
		// 	return fmt.Errorf("failed to marshal JSON: %w", err)
		// }
		// fmt.Println(string(prettyJSON))
		displayStructure, logicalStructure := getGenericTableStructures(attributes)

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{displayStructure.Id, displayStructure.PrimaryName, displayStructure.CreatedOn, displayStructure.State, displayStructure.Owner})
		fmt.Println(logicalStructure)

		for _, record := range records {
			t.AppendRow(table.Row{
				record[logicalStructure.Id],
				record[logicalStructure.PrimaryName],
				record[logicalStructure.CreatedOn],
				record[logicalStructure.State],
				record[logicalStructure.Owner],
			})
		}

		t.Render()
	}

	return nil
}

func getGenericTableStructures(attributes []models.EntityAttribute) (models.GenericTableStructure, models.GenericTableStructure) {
	displayStructure := models.GenericTableStructure{Id: "ID", CreatedOn: "Created On", State: "State", Owner: "Owner"}
	logicalStructure := models.GenericTableStructure{CreatedOn: "createdon", State: "state", Owner: "owner"}

	for _, attr := range attributes {
		if attr.IsPrimaryName {
			displayStructure.PrimaryName = attr.DisplayName
			logicalStructure.PrimaryName = attr.LogicalName
		}
		// TODO: check for AttributeOf (null)
		if attr.IsPrimaryId {
			logicalStructure.PrimaryName = attr.LogicalName
		}
	}

	return displayStructure, logicalStructure
}

func parseOptions(args []string) (models.ListOptions, error) {
	options := models.QueryOptions{}
	format := models.OutputFormatTable
	var rawFormat string

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
		case "--output":
			if err = parseStringFlag(&rawFormat, args[i], &i); err != nil {
				break
			}
			switch strings.ToLower(rawFormat) {
			case "json":
				format = models.OutputFormatJSON
			case "table":
				format = models.OutputFormatTable
			default:
				return models.ListOptions{}, fmt.Errorf("invalid output format %q (allowed: json, table)", rawFormat)
			}
		default:
			return models.ListOptions{}, fmt.Errorf("unknown option: %s", args[i])
		}

		if err != nil {
			return models.ListOptions{}, err
		}
	}

	return models.ListOptions{Query: options, Output: format}, nil
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
