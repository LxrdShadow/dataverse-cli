package commands

import "dvc/client"

var MetadataCommand = &Command{
	Name:        "metadata",
	Description: "Get the metadata of an entity",
	Run:         runMetadataCommand,
}

func runMetadataCommand(client *client.DataverseClient, args []string) error {

	return nil
}
