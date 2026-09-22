package commands

import "dvc/client"

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

	return nil
}

func metadataUsage() string {
	return `Usage: dvc metadata <entity>`
}
