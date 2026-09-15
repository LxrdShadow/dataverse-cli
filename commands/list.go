package commands

import (
	"dvc/client"
	"dvc/models"
	"encoding/json"
	"fmt"
)

var ListCommand = &Command{
	Name:        "list",
	Description: "Retrieve records from a Dataverse table",
	Run:         listRecords,
}

func listRecords(client *client.DataverseClient, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: dvc list <table_logical_name>")
	}

	table, err := client.GetTable(args[0])
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	records, err := client.ListRecords(table.EntitySetName, models.QueryOptions{})
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
