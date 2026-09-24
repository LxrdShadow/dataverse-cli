package commands

import (
	"dvc/models"
)

func filterAttributes(
	attributes []models.EntityAttribute,
	requiredLevel models.RequirementLevel,
) []models.EntityAttribute {
	if requiredLevel == "" {
		return attributes
	}

	var filtered []models.EntityAttribute

	for _, attribute := range attributes {
		if attribute.RequiredLevel == requiredLevel {
			filtered = append(filtered, attribute)
		}
	}

	return filtered
}
