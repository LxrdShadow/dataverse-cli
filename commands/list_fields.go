package commands

import (
	"dvc/models"
	"strings"
)

type TableField struct {
	DisplayName string
	QueryName   string
	RecordName  string
}

func getMetadataFieldName(fieldName string) string {
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

func getQueryFieldName(attr models.EntityAttribute) string {
	switch attr.AttributeType {
	case models.AttributeTypeOwner:
		return "_" + attr.LogicalName + "_value"

	default:
		return attr.LogicalName
	}
}

func getSelectFields(attributes []models.EntityAttribute, selectedFields string) []TableField {
	var fields []TableField

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

		metadataName := getMetadataFieldName(fieldName)
		attr, found := availableAttributes[metadataName]
		if !found {
			continue
		}

		fields = append(fields, TableField{
			DisplayName: attr.DisplayName,
			QueryName:   getQueryFieldName(attr),
			RecordName:  getRecordFieldName(attr),
		})
	}
	return fields
}

func getRecordFieldName(attr models.EntityAttribute) string {
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

func getDefaultTableFields(
	attributes []models.EntityAttribute,
) ([]TableField, []TableField) {
	var fields []TableField

	for _, attr := range attributes {
		switch {
		case attr.IsPrimaryId && !attr.IsLogical:
			fields = append(fields, TableField{
				DisplayName: "ID",
				QueryName:   attr.LogicalName,
				RecordName:  attr.LogicalName,
			})

		case attr.IsPrimaryName:
			fields = append(fields, TableField{
				DisplayName: attr.DisplayName,
				QueryName:   attr.LogicalName,
				RecordName:  attr.LogicalName,
			})

		case attr.LogicalName == "createdon":
			fields = append(fields, TableField{
				DisplayName: "Created On",
				QueryName:   attr.LogicalName,
				RecordName:  getRecordFieldName(attr),
			})

		case attr.LogicalName == "statecode":
			fields = append(fields, TableField{
				DisplayName: "State",
				QueryName:   attr.LogicalName,
				RecordName:  getRecordFieldName(attr),
			})

		case attr.LogicalName == "ownerid":
			fields = append(fields, TableField{
				DisplayName: "Owner",
				QueryName:   getQueryFieldName(attr),
				RecordName:  getRecordFieldName(attr),
			})
		}
	}

	return fields, fields
}
