package models

import "encoding/json/v2"

type Entity struct {
	LogicalName          string
	DisplayName          string
	EntitySetName        string
	IsCustomEntity       bool
	IsManaged            bool
	PrimaryIdAttribute   string
	PrimaryNameAttribute string
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	var raw struct {
		LogicalName          string
		EntitySetName        string
		DisplayName          DisplayName
		IsCustomEntity       bool
		IsManaged            bool
		PrimaryIdAttribute   string
		PrimaryNameAttribute string
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	e.LogicalName = raw.LogicalName
	e.EntitySetName = raw.EntitySetName
	e.DisplayName = raw.DisplayName.UserLocalizedLabel.Label
	e.IsCustomEntity = raw.IsCustomEntity
	e.IsManaged = raw.IsManaged
	e.PrimaryIdAttribute = raw.PrimaryIdAttribute
	e.PrimaryNameAttribute = raw.PrimaryNameAttribute
	return nil
}
