package models

import "encoding/json/v2"

type AttributeType string

const (
	AttributeTypePicklist AttributeType = "Picklist"
	AttributeTypeState    AttributeType = "State"
	AttributeTypeOwner    AttributeType = "Owner"
	AttributeTypeDateTime AttributeType = "DateTime"
)

type RequiredLevel struct {
	Value string
}

type EntityAttribute struct {
	LogicalName   string
	DisplayName   string
	IsPrimaryName bool
	IsPrimaryId   bool
	IsLogical     bool
	RequiredLevel string
	AttributeType AttributeType
}

func (a *EntityAttribute) UnmarshalJSON(data []byte) error {
	var raw struct {
		LogicalName   string
		DisplayName   DisplayName
		IsPrimaryName bool
		IsPrimaryId   bool
		IsLogical     bool
		RequiredLevel RequiredLevel
		AttributeType AttributeType
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	a.LogicalName = raw.LogicalName
	a.DisplayName = raw.DisplayName.UserLocalizedLabel.Label
	a.IsPrimaryName = raw.IsPrimaryName
	a.IsPrimaryId = raw.IsPrimaryId
	a.IsLogical = raw.IsLogical
	a.RequiredLevel = raw.RequiredLevel.Value
	a.AttributeType = raw.AttributeType
	return nil
}
