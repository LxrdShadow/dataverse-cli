package commands

import "dvc/client"

type Command struct {
	Name        string
	Description string
	Run         func(*client.DataverseClient, []string) error
}
