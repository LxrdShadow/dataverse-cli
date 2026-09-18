package commands

import (
	"dvc/client"
	"dvc/models"
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
		return err
	}

	entities, err := client.ListTables(filters.Scope, filters.Management)
	if err != nil {
		return err
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
	filters := TableFilters{
		Scope:      models.TableScopeAll,
		Management: models.TableManagementAll,
	}

	for _, arg := range args {
		switch arg {
		case "--custom":
			if filters.Scope == models.TableScopeSystem {
				return TableFilters{}, fmt.Errorf("scope cannot be both system and custom")
			}
			filters.Scope = models.TableScopeCustom

		case "--system":
			if filters.Scope == models.TableScopeCustom {
				return TableFilters{}, fmt.Errorf("scope cannot be both system and custom")
			}
			filters.Scope = models.TableScopeSystem

		case "--managed":
			if filters.Management == models.TableManagementUnmanaged {
				return TableFilters{}, fmt.Errorf("management cannot be both unmanaged and managed")
			}
			filters.Management = models.TableManagementManaged

		case "--unmanaged":
			if filters.Management == models.TableManagementManaged {
				return TableFilters{}, fmt.Errorf("management cannot be both unmanaged and managed")
			}
			filters.Management = models.TableManagementUnmanaged

		default:
			return TableFilters{}, fmt.Errorf("unknown option: %s", arg)
		}
	}

	return filters, nil
}
