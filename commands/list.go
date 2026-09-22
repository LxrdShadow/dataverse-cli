package commands

import (
	"dvc/client"
	"dvc/models"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

var ListCommand = &Command{
	Name:        "list",
	Description: "Retrieve records from a Dataverse table",
	Run:         listRecords,
}

func listRecords(client *client.DataverseClient, args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Println(listUsage())
		return nil
	}

	logicalName := args[0]
	optionArgs := args[1:]
	listOptions, err := parseListOptions(optionArgs)
	if err != nil {
		return err
	}

	entity, err := client.GetTable(logicalName, true)
	if err != nil {
		return err
	}

	if listOptions.Query.Select == "" && listOptions.Output == models.OutputFormatTable {
		defaultFields := getDefaultTableFields(entity.Attributes)

		queryFields := make([]string, 0, len(defaultFields))
		for _, field := range defaultFields {
			queryFields = append(queryFields, field.QueryName)
		}

		listOptions.Query.Select = strings.Join(queryFields, ",")
	}

	records, err := client.ListRecords(entity.EntitySetName, listOptions)
	if err != nil {
		return fmt.Errorf("failed to list records: %w", err)
	}

	switch listOptions.Output {
	case models.OutputFormatJSON:
		return renderJSON(records)
	case models.OutputFormatTable:
		err = renderTable(entity.Attributes, listOptions.Query.Select, records)
		if err != nil {
			return fmt.Errorf("failed to render table: %w", err)
		}
		fmt.Printf("Retrieved %d records\n", len(records))
	default:
		return fmt.Errorf("unsupported output format: %s", listOptions.Output)
	}

	return nil
}

func renderTable(
	attributes []models.EntityAttribute,
	selectedFields string,
	records []models.Record,
) error {
	fields := getSelectFields(attributes, selectedFields)
	if len(fields) == 0 {
		return fmt.Errorf("none of the selected columns exist on the table")
	}

	headerRow := table.Row{}
	for _, field := range fields {
		headerRow = append(headerRow, field.DisplayName)
	}

	rows := make([]table.Row, 0, len(records))
	for _, record := range records {
		row := table.Row{}

		for _, field := range fields {
			value := record[field.RecordName]
			row = append(row, value)
		}

		rows = append(rows, row)
	}

	printTable(headerRow, rows)
	return nil
}

func listUsage() string {
	return `Usage: dvc list <table> [options]

Options:
  -h, --help				Show this help message
  --select <columns>			Select columns to retrieve and display (default: common columns)
  --filter <condition>			Filter rows based on an OData condition
  --order-by <column>			Order rows by a specific column
  --expand <columns>			Expand related columns
  --top <count>				Limit the total number of rows returned
  --max-page-size <count>		Limit the number of rows per query page
  -o --output <format>			Output format (default: table, possible values: json, table)`
}
