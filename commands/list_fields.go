package commands

import (
	"dvc/models"
	"strings"
)

type DisplayField struct {
	DisplayName string
	QueryName   string
	RecordName  string
}

func logicalNameFromQueryName(fieldName string) string {
	switch {
	case strings.HasPrefix(fieldName, "_") &&
		strings.HasSuffix(fieldName, "_value"):
		return strings.TrimSuffix(
			strings.TrimPrefix(fieldName, "_"),
			"_value",
		)

	default:
		return fieldName
	}
}

func queryNameForAttribute(attr models.EntityAttribute) string {
	switch attr.AttributeType {
	case models.AttributeTypeOwner:
		return "_" + attr.LogicalName + "_value"

	default:
		return attr.LogicalName
	}
}

func getSelectFields(attributes []models.EntityAttribute, selectedFields string) []DisplayField {
	var fields []DisplayField

	availableAttributes := make(map[string]models.EntityAttribute)
	for _, attr := range attributes {
		availableAttributes[attr.LogicalName] = attr
	}

	selected := strings.SplitSeq(selectedFields, ",")
	for fieldName := range selected {
		fieldName = strings.TrimSpace(fieldName)
		if fieldName == "" {
			continue
		}

		metadataName := logicalNameFromQueryName(fieldName)
		attr, found := availableAttributes[metadataName]
		if !found {
			continue
		}

		fields = append(fields, DisplayField{
			DisplayName: attr.DisplayName,
			QueryName:   queryNameForAttribute(attr),
			RecordName:  recordValueNameForAttribute(attr),
		})
	}
	return fields
}

func recordValueNameForAttribute(attr models.EntityAttribute) string {
	const formattedSuffix = "@OData.Community.Display.V1.FormattedValue"

	switch attr.AttributeType {
	case models.AttributeTypeOwner:
		return "_ownerid_value" + formattedSuffix

	case models.AttributeTypePicklist,
		models.AttributeTypeState,
		models.AttributeTypeDateTime:
		return attr.LogicalName + formattedSuffix

	default:
		return attr.LogicalName
	}
}

func defaultDisplayFields(
	attributes []models.EntityAttribute,
) []DisplayField {
	var fields []DisplayField

	for _, attr := range attributes {
		switch {
		case attr.IsPrimaryId && !attr.IsLogical:
			fields = append(fields, DisplayField{
				DisplayName: "ID",
				QueryName:   attr.LogicalName,
				RecordName:  attr.LogicalName,
			})

		case attr.IsPrimaryName:
			fields = append(fields, DisplayField{
				DisplayName: attr.DisplayName,
				QueryName:   attr.LogicalName,
				RecordName:  attr.LogicalName,
			})

		case attr.LogicalName == "createdon":
			fields = append(fields, DisplayField{
				DisplayName: "Created On",
				QueryName:   attr.LogicalName,
				RecordName:  recordValueNameForAttribute(attr),
			})

		case attr.LogicalName == "statecode":
			fields = append(fields, DisplayField{
				DisplayName: "State",
				QueryName:   attr.LogicalName,
				RecordName:  recordValueNameForAttribute(attr),
			})

		case attr.LogicalName == "ownerid":
			fields = append(fields, DisplayField{
				DisplayName: "Owner",
				QueryName:   queryNameForAttribute(attr),
				RecordName:  recordValueNameForAttribute(attr),
			})
		}
	}

	return fields
}
