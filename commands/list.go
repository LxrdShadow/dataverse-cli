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

	if listOptions.Query.Select == "" && listOptions.Output == models.OutputFormatTable {
		_, defaultFields := getDefaultTableFields(entity.Attributes)

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

	if len(records) == 0 {
		fmt.Println("no records found")
		return nil
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

func renderJSON(records []models.Record) error {
	prettyJSON, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(prettyJSON))

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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	headerRow := table.Row{}

	for _, field := range fields {
		headerRow = append(headerRow, field.DisplayName)
	}

	t.AppendHeader(headerRow)

	for _, record := range records {
		row := table.Row{}

		for _, field := range fields {
			value := record[field.RecordName]
			row = append(row, value)
		}

		t.AppendRow(row)
	}

	t.Render()

	return nil
}

func getMetadataFieldName(fieldName string) string {
	switch {
	case strings.HasPrefix(fieldName, "_") &&
		strings.HasSuffix(fieldName, "_value"):
		return strings.TrimSuffix(
			strings.TrimPrefix(fieldName, "_"),
			"_value",
		)

	default:
		return fieldName
	}
}

func getQueryFieldName(attr models.EntityAttribute) string {
	switch attr.AttributeType {
	case models.AttributeTypeOwner:
		return "_" + attr.LogicalName + "_value"

	default:
		return attr.LogicalName
	}
}

func getSelectFields(attributes []models.EntityAttribute, selectedFields string) []models.TableField {
	var fields []models.TableField

	availableAttributes := make(map[string]models.EntityAttribute)
	for _, attr := range attributes {
		availableAttributes[attr.LogicalName] = attr
	}

	selected := strings.SplitSeq(selectedFields, ",")
	for fieldName := range selected {
		fieldName = strings.TrimSpace(fieldName)
		if fieldName == "" {
			continue
		}

		metadataName := getMetadataFieldName(fieldName)
		attr, found := availableAttributes[metadataName]
		if !found {
			continue
		}

		fields = append(fields, models.TableField{
			DisplayName: attr.DisplayName,
			QueryName:   getQueryFieldName(attr),
			RecordName:  getRecordFieldName(attr),
		})
	}
	return fields
}

func getRecordFieldName(attr models.EntityAttribute) string {
	const formattedSuffix = "@OData.Community.Display.V1.FormattedValue"

	switch attr.AttributeType {
	case models.AttributeTypeOwner:
		return "_ownerid_value" + formattedSuffix

	case models.AttributeTypePicklist,
		models.AttributeTypeState,
		models.AttributeTypeDateTime:
		return attr.LogicalName + formattedSuffix

	default:
		return attr.LogicalName
	}
}

func getDefaultTableFields(
	attributes []models.EntityAttribute,
) ([]models.TableField, []models.TableField) {
	var fields []models.TableField

	for _, attr := range attributes {
		switch {
		case attr.IsPrimaryId && !attr.IsLogical:
			fields = append(fields, models.TableField{
				DisplayName: "ID",
				QueryName:   attr.LogicalName,
				RecordName:  attr.LogicalName,
			})

		case attr.IsPrimaryName:
			fields = append(fields, models.TableField{
				DisplayName: attr.DisplayName,
				QueryName:   attr.LogicalName,
				RecordName:  attr.LogicalName,
			})

		case attr.LogicalName == "createdon":
			fields = append(fields, models.TableField{
				DisplayName: "Created On",
				QueryName:   attr.LogicalName,
				RecordName:  getRecordFieldName(attr),
			})

		case attr.LogicalName == "statecode":
			fields = append(fields, models.TableField{
				DisplayName: "State",
				QueryName:   attr.LogicalName,
				RecordName:  getRecordFieldName(attr),
			})

		case attr.LogicalName == "ownerid":
			fields = append(fields, models.TableField{
				DisplayName: "Owner",
				QueryName:   getQueryFieldName(attr),
				RecordName:  getRecordFieldName(attr),
			})
		}
	}

	return fields, fields
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
			err = parseStringFlag(&options.Select, args[i], &i)
		case "--filter":
			err = parseStringFlag(&options.Filter, args[i], &i)
		case "--order-by":
			err = parseStringFlag(&options.OrderBy, args[i], &i)
		case "--expand":
			err = parseStringFlag(&options.Expand, args[i], &i)
		case "--top":
			err = parseIntFlag(&options.Top, args[i], &i)
		case "--max-page-size":
			err = parseIntFlag(&options.MaxPageSize, args[i], &i)
		case "--output", "-o":
			if err = parseStringFlag(&rawFormat, args[i], &i); err != nil {
				break
			}
			switch strings.ToLower(rawFormat) {
			case string(models.OutputFormatJSON):
				format = models.OutputFormatJSON
			case string(models.OutputFormatTable):
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
  -h, --help				Show this help message
  --select <columns>			Select columns to retrieve and display (default: common columns)
  --filter <condition>			Filter rows based on an OData condition
  --order-by <column>			Order rows by a specific column
  --expand <columns>			Expand related columns
  --top <count>				Limit the total number of rows returned
  --max-page-size <count>		Limit the number of rows per query page
  -o --output <format>			Output format (default: table, possible values: json, table)`
}
