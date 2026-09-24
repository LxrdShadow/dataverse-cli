package client

import (
	"fmt"
	"strings"
)

func odataStringLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func entityDefinitionPath(logicalName string) string {
	return fmt.Sprintf(
		"EntityDefinitions(LogicalName=%s)",
		odataStringLiteral(logicalName),
	)
}
