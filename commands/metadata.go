package commands

import (
	"dvc/client"
	"dvc/models"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

var MetadataCommand = &Command{
	Name:        "metadata",
	Description: "Get the metadata of an entity",
	Run:         runMetadataCommand,
}

func runMetadataCommand(client *client.DataverseClient, args []string) error {
	err := validateArgs(args, 1, metadataUsage())
	if err != nil {
		return err
	}

	options, err := parseEntityMetadataOptions(args[1:])
	if err != nil {
		return err
	}

	if options.Attributes {
		attributes, err := client.ListEntityAttributes(args[0])
		if err != nil {
			return err
		}

		switch options.Output {
		case models.OutputFormatJSON:
			return renderJSON(attributes)
		case models.OutputFormatTable:
			if err := renderAttributesMetadataTable(attributes); err != nil {
				return fmt.Errorf("failed to render table: %w", err)
			}
		default:
			return fmt.Errorf("unsupported output format: %s", options.Output)
		}
	} else {
		entity, err := client.GetEntity(args[0])
		if err != nil {
			return err
		}

		switch options.Output {
		case models.OutputFormatJSON:
			return renderJSON(entity)
		case models.OutputFormatTable:
			if err := renderEntityMetadataTable(entity); err != nil {
				return fmt.Errorf("failed to render table: %w", err)
			}
		default:
			return fmt.Errorf("unsupported output format: %s", options.Output)
		}
	}

	return nil
}

func renderEntityMetadataTable(entity *models.Entity) error {
	headerRow := table.Row{"Entity Metadata", "Value"}
	rows := []table.Row{
		{"Logical Name", entity.LogicalName},
		{"Display Name", entity.DisplayName},
		{"Entity Set Name", entity.EntitySetName},
		{"Is Custom Entity", fmt.Sprintf("%t", entity.IsCustomEntity)},
		{"Is Managed", fmt.Sprintf("%t", entity.IsManaged)},
		{"Primary Id Attribute", entity.PrimaryIdAttribute},
		{"Primary Name Attribute", entity.PrimaryNameAttribute},
	}

	printTable(headerRow, rows)
	return nil
}

func renderAttributesMetadataTable(attributes []models.EntityAttribute) error {
	headerRow := table.Row{"Logical Name", "Display Name", "Type", "Required"}
	rows := make([]table.Row, 0, len(attributes))
	for _, attr := range attributes {
		rows = append(rows, table.Row{
			attr.LogicalName,
			attr.DisplayName,
			attr.AttributeType,
			getRequiredLevel(attr.RequiredLevel),
		})
	}

	printTable(headerRow, rows)
	return nil
}

func getRequiredLevel(requiredLevel string) string {
	var requiredLevelMap = map[string]string{
		"None":                "None",
		"SystemRequired":      "System Required",
		"ApplicationRequired": "Application Required",
		"Recommended":         "Recommended",
	}

	if level, ok := requiredLevelMap[requiredLevel]; ok {
		return level
	}
	return requiredLevel
}

func metadataUsage() string {
	return `Usage: dvc metadata <entity> [options]

Options:
  -h, --help   			Show this help message
  --attributes			Display the metadata of the attributes of the entity
  -o, --output <format>  	Output format (json, table) (default: table)`
}
