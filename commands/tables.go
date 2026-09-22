package commands

import (
	"dvc/client"
	"dvc/models"
	"flag"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

type TableFilters struct {
	Scope      models.EntityScope
	Management models.EntityManagement
}

var TablesCommand = &Command{
	Name:        "tables",
	Description: "Retrieves a list of the tables in the current Dataverse environment.",
	Run:         runTablesCommand,
}

func runTablesCommand(client *client.DataverseClient, args []string) error {
	filters, err := parseTableFilters(args)
	if err != nil {
		return err
	}

	entities, err := client.ListEntities(filters.Scope, filters.Management)
	if err != nil {
		return fmt.Errorf("failed to fetch table list: %w", err)
	}

	if len(entities) == 0 {
		fmt.Println("No tables matched the specified filters")
		return nil
	}

	headerRow := table.Row{"Logical Name", "Display Name", "Is Custom", "Is Managed"}

	rows := make([]table.Row, 0, len(entities))
	for _, entity := range entities {
		rows = append(rows, table.Row{
			entity.LogicalName,
			entity.DisplayName,
			entity.IsCustomEntity,
			entity.IsManaged,
		})
	}

	printTable(headerRow, rows)

	fmt.Printf("Retrieved %d tables.\n", len(entities))

	return nil
}

func parseTableFilters(args []string) (TableFilters, error) {
	fs := flag.NewFlagSet("tables", flag.ContinueOnError)
	var custom, system, managed, unmanaged bool
	fs.BoolVar(&custom, "custom", false, "")
	fs.BoolVar(&system, "system", false, "")
	fs.BoolVar(&managed, "managed", false, "")
	fs.BoolVar(&unmanaged, "unmanaged", false, "")
	fs.Usage = func() { fmt.Println(tablesUsage()) }

	if err := fs.Parse(args); err != nil {
		return TableFilters{}, err
	}
	if custom && system {
		return TableFilters{}, fmt.Errorf("scope cannot be both system and custom")
	}
	if managed && unmanaged {
		return TableFilters{}, fmt.Errorf("management cannot be both managed and unmanaged")
	}

	filters := TableFilters{Scope: models.EntityScopeAll, Management: models.EntityManagementAll}
	if custom {
		filters.Scope = models.EntityScopeCustom
	} else if system {
		filters.Scope = models.EntityScopeSystem
	}
	if managed {
		filters.Management = models.EntityManagementManaged
	} else if unmanaged {
		filters.Management = models.EntityManagementUnmanaged
	}

	return filters, nil
}

func tablesUsage() string {
	return `Usage:
  dvc tables [options]

Options:
  -custom    List custom tables
  -system    List system tables
  -managed   List managed tables
  -unmanaged List unmanaged tables
  -h, --help Show this help message`
}
