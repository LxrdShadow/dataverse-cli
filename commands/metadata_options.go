package commands

import (
	"dvc/models"
	"flag"
	"fmt"
	"strings"
)

func parseEntityMetadataOptions(args []string) (models.EntityMetadataOptions, error) {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)
	var attributes, relationships bool
	var rawFormat, rawRequirementLevel string

	fs.BoolVar(&attributes, "attributes", false, "Display the metadata of the attributes of the entity")
	fs.BoolVar(&relationships, "relationships", false, "Display the metadata of the relationships of the entity")
	fs.StringVar(&rawFormat, "output", "table", "Output format (table, json)")
	fs.StringVar(&rawFormat, "o", "table", "Output format (table, json) shorthand")
	fs.StringVar(&rawRequirementLevel, "required-level", "", "Requirement level (none, systemrequired, applicationrequired, recommended)")
	fs.Usage = func() { fmt.Println(metadataUsage()) }

	if err := fs.Parse(args); err != nil {
		return models.EntityMetadataOptions{}, err
	}

	if attributes && relationships {
		return models.EntityMetadataOptions{}, fmt.Errorf("attributes and relationships cannot be displayed at the same time")
	}

	if rawRequirementLevel != "" && !attributes {
		return models.EntityMetadataOptions{}, fmt.Errorf("requirement-level can only be specified with attributes")
	}

	format, err := parseOutputFormat(rawFormat)
	if err != nil {
		return models.EntityMetadataOptions{}, err
	}

	requirementLevel, err := parseRequirementLevel(rawRequirementLevel)
	if err != nil {
		return models.EntityMetadataOptions{}, err
	}

	return models.EntityMetadataOptions{
		Output:              format,
		Attributes:          attributes,
		Relationships:       relationships,
		RequirementLevel:    requirementLevel,
		RequirementLevelSet: rawRequirementLevel != "",
	}, nil
}

func parseOutputFormat(rawFormat string) (models.OutputFormat, error) {
	switch strings.ToLower(rawFormat) {
	case string(models.OutputFormatJSON):
		return models.OutputFormatJSON, nil
	case string(models.OutputFormatTable):
		return models.OutputFormatTable, nil
	default:
		return models.OutputFormatTable, fmt.Errorf("invalid output format %q (allowed: json, table)", rawFormat)
	}
}

func parseRequirementLevel(rawRequirementLevel string) (models.RequirementLevel, error) {
	var requirementLevelMap = map[string]models.RequirementLevel{
		"":            models.RequirementLevelNone,
		"none":        models.RequirementLevelNone,
		"system":      models.RequirementLevelSystemRequired,
		"application": models.RequirementLevelApplicationRequired,
		"recommended": models.RequirementLevelRecommended,
	}

	if parsed, ok := requirementLevelMap[strings.ToLower(rawRequirementLevel)]; ok {
		return parsed, nil
	}

	parsed := models.RequirementLevel(rawRequirementLevel)
	switch parsed {
	case models.RequirementLevelNone,
		models.RequirementLevelSystemRequired,
		models.RequirementLevelApplicationRequired,
		models.RequirementLevelRecommended:
		return parsed, nil
	default:
		return models.RequirementLevelNone, fmt.Errorf("invalid requirement level %q. Allowed: None, SystemRequired, ApplicationRequired, Recommended (or none, system, application, recommended)", rawRequirementLevel)
	}
}
