package commands

import (
	"dvc/client"
	"dvc/models"
	"errors"
	"flag"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

type TableFilters struct {
	Scope      models.TableScope
	Management models.TableManagement
}

var TablesCommand = &Command{
	Name:        "tables",
	Description: "Retrieves a list of the tables in the current Dataverse environment.",
	Run:         tables,
}

func tables(client *client.DataverseClient, args []string) error {
	filters, err := parseTableFilters(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	entities, err := client.ListTables(filters.Scope, filters.Management)
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
			entity.IsCustom,
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

	filters := TableFilters{Scope: models.TableScopeAll, Management: models.TableManagementAll}
	if custom {
		filters.Scope = models.TableScopeCustom
	} else if system {
		filters.Scope = models.TableScopeSystem
	}
	if managed {
		filters.Management = models.TableManagementManaged
	} else if unmanaged {
		filters.Management = models.TableManagementUnmanaged
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
