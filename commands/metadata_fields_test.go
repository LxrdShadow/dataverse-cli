package commands

import (
	"reflect"
	"testing"

	"dvc/models"
)

func TestFilterAttributes(t *testing.T) {
	attributes := []models.EntityAttribute{
		{LogicalName: "optional", RequiredLevel: models.RequirementLevelNone},
		{LogicalName: "system", RequiredLevel: models.RequirementLevelSystemRequired},
		{LogicalName: "recommended", RequiredLevel: models.RequirementLevelRecommended},
	}

	tests := []struct {
		name          string
		requiredLevel models.RequirementLevel
		want          []models.EntityAttribute
	}{
		{name: "no filter", requiredLevel: "", want: attributes},
		{name: "optional", requiredLevel: models.RequirementLevelNone, want: attributes[:1]},
		{name: "system required", requiredLevel: models.RequirementLevelSystemRequired, want: attributes[1:2]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterAttributes(attributes, tt.requiredLevel)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterAttributes(_, %q) = %+v, want %+v", tt.requiredLevel, got, tt.want)
			}
		})
	}
}
