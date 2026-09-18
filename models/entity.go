package models

import "encoding/json/v2"

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
