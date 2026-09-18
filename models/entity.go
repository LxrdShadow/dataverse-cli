package models

import "encoding/json/v2"

type DisplayName struct {
	UserLocalizedLabel LocalizedLabel `json:"UserLocalizedLabel"`
}

type LocalizedLabel struct {
	Label string `json:"Label"`
}

type Entity struct {
	LogicalName   string
	DisplayName   string
	EntitySetName string
	IsCustom      bool
	IsManaged     bool
	Attributes    []EntityAttribute
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	var raw struct {
		LogicalName    string
		EntitySetName  string
		DisplayName    DisplayName
		IsCustomEntity bool
		IsManaged      bool
		Attributes     []EntityAttribute
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	e.LogicalName = raw.LogicalName
	e.EntitySetName = raw.EntitySetName
	e.DisplayName = raw.DisplayName.UserLocalizedLabel.Label
	e.IsCustom = raw.IsCustomEntity
	e.IsManaged = raw.IsManaged
	e.Attributes = raw.Attributes
	return nil
}

type AttributeType string

const (
	AttributeTypePicklist AttributeType = "Picklist"
	AttributeTypeState    AttributeType = "State"
	AttributeTypeOwner    AttributeType = "Owner"
	AttributeTypeDateTime AttributeType = "DateTime"
)

type EntityAttribute struct {
	LogicalName   string
	DisplayName   string
	IsPrimaryName bool
	IsPrimaryId   bool
	IsLogical     bool
	AttributeType AttributeType
}

func (a *EntityAttribute) UnmarshalJSON(data []byte) error {
	var raw struct {
		LogicalName   string
		DisplayName   DisplayName
		IsPrimaryName bool
		IsPrimaryId   bool
		IsLogical     bool
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
	a.AttributeType = raw.AttributeType
	return nil
}
