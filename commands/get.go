package commands

import (
	"dvc/client"
	"dvc/models"
	"dvc/utils"
	"fmt"
	"uuid"

	"github.com/jedib0t/go-pretty/v6/table"
)

var GetCommand = &Command{
	Name:        "get",
	Description: "Get a record by ID",
	Run:         getRecord,
}

func getRecord(client *client.DataverseClient, args []string) error {
	if len(args) < 2 || args[0] == "-h" || args[0] == "--help" {
		fmt.Println(getUsage())
		return nil
	}

	logicalName := args[0]
	recordID := args[1]

	_, err := uuid.Parse(recordID)
	if err != nil {
		return fmt.Errorf("invalid record ID: %w", err)
	}

	options, err := parseGetOptions(args[2:])
	if err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	entity, err := client.GetTable(logicalName, true)
	if err != nil {
		return err
	}

	record, err := client.GetRecord(entity.EntitySetName, recordID, options)
	if err != nil {
		return fmt.Errorf("failed to list records: %w", err)
	}

	if record == nil {
		fmt.Println("record not found")
		return nil
	}

	switch options.Output {
	case models.OutputFormatJSON:
		return renderJSON(record)
	case models.OutputFormatTable:
		err := renderFieldValueTable(entity.Attributes, options.Select, record)
		if err != nil {
			return fmt.Errorf("failed to render table: %w", err)
		}
	default:
		return fmt.Errorf("unsupported output format: %s", options.Output)
	}

	return nil
}

func renderFieldValueTable(attributes []models.EntityAttribute, selectedFields string, record models.Record) error {
	headerRow := table.Row{}
	rows := make([]table.Row, 0, len(record))

	fields := getSelectFields(attributes, selectedFields)
	headerRow = table.Row{"Field", "Display Name", "Value"}
	for _, field := range fields {
		row := table.Row{field.QueryName, field.DisplayName, utils.ValueOrEmpty(record, field.RecordName)}
		rows = append(rows, row)
	}

	printTable(headerRow, rows)
	return nil
}

func getUsage() string {
	return `Usage: dvc get <table> <recordID>

Options:
	--select <fields>  Comma-separated list of fields to select
	--output, -o <format>  Output format (json, table) (default: table)`
}
