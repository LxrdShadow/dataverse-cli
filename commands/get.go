package commands

import (
	"dvc/client"
	"dvc/models"
	"fmt"
	"uuid"

	"github.com/jedib0t/go-pretty/v6/table"
)

var GetCommand = &Command{
	Name:        "get",
	Description: "Get a record by ID",
	Run:         runGetCommand,
}

func runGetCommand(client *client.DataverseClient, args []string) error {
	err := validateArgs(args, 2, getUsage())
	if err != nil {
		return err
	}

	logicalName := args[0]
	recordID := args[1]

	_, err = uuid.Parse(recordID)
	if err != nil {
		return fmt.Errorf("invalid record ID: %w", err)
	}

	options, err := parseGetRecordOptions(args[2:])
	if err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	entity, err := client.GetEntity(logicalName, true)
	if err != nil {
		return err
	}

	record, err := client.GetRecord(entity.EntitySetName, recordID, options)
	if err != nil {
		return fmt.Errorf("failed to get record: %w", err)
	}

	switch options.Output {
	case models.OutputFormatJSON:
		return renderJSON(record)
	case models.OutputFormatTable:
		err := renderFieldValueTable(entity.Attributes, options, record)
		if err != nil {
			return fmt.Errorf("failed to render table: %w", err)
		}
	default:
		return fmt.Errorf("unsupported output format: %s", options.Output)
	}

	return nil
}

func renderFieldValueTable(attributes []models.EntityAttribute, options models.GetRecordOptions, record models.Record) error {
	headerRow := table.Row{}
	rows := make([]table.Row, 0, len(record))

	fields := getSelectFields(attributes, options.Select)
	if len(fields) == 0 && options.Select == "" {
		for _, attr := range attributes {
			fields = append(fields, DisplayField{
				DisplayName: attr.DisplayName,
				QueryName:   queryNameForAttribute(attr),
				RecordName:  recordValueNameForAttribute(attr),
			})
		}
	} else if len(fields) == 0 {
		return fmt.Errorf("none of the selected columns exist on the table")
	}

	headerRow = table.Row{"Field", "Display Name", "Value"}
	for _, field := range fields {
		if value := record[field.RecordName]; value == nil && !options.DisplayEmpty {
			continue
		}
		row := table.Row{field.QueryName, field.DisplayName, valueOrEmpty(record, field.RecordName)}
		rows = append(rows, row)
	}

	printTable(headerRow, rows)
	return nil
}

func getUsage() string {
	return `Usage: dvc get <table> <recordID>

Options:
  -h, --help			Show this help message
  --select <fields>  		Comma-separated list of fields to select
  --output, -o <format>  	Output format (json, table) (default: table)
  --display-empty  		Display empty values (default: false)
`
}
