package commands

var CommandRegistry = map[string]*Command{
	WhoAmICommand.Name:   WhoAmICommand,
	TablesCommand.Name:   TablesCommand,
	ListCommand.Name:     ListCommand,
	GetCommand.Name:      GetCommand,
	MetadataCommand.Name: MetadataCommand,
}
