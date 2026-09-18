package commands

import (
	"encoding/json"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func printTable(headers table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(headers)
	for _, row := range rows {
		t.AppendRow(row)
	}
	t.Render()
}

func renderJSON(data any) error {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(prettyJSON)
	return err
}
