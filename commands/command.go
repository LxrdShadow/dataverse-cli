package commands

import "dvc/client"

type Command struct {
	Name        string
	Description string
	Run         func(*client.DataverseClient, []string) error
}

func (c *Command) Execute(client *client.DataverseClient, args []string) error {
	return c.Run(client, args)
}
