package commands

import (
	"reflect"
	"testing"

	"dvc/models"
)

func TestGetMetadataFieldName(t *testing.T) {
	tests := []struct {
		name  string
		field string
		want  string
	}{
		{"lookup value field", "_ownerid_value", "ownerid"},
		{"plain field", "name", "name"},
		{"prefix without suffix", "_ownerid", "_ownerid"},
		{"suffix without prefix", "ownerid_value", "ownerid_value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := logicalNameFromQueryName(tt.field)
			if got != tt.want {
				t.Errorf("getMetadataFieldName(%q) = %q, want %q", tt.field, got, tt.want)
			}
		})
	}
}

func TestGetQueryFieldName(t *testing.T) {
	tests := []struct {
		name string
		attr models.EntityAttribute
		want string
	}{
		{
			name: "owner attribute",
			attr: models.EntityAttribute{LogicalName: "ownerid", AttributeType: models.AttributeTypeOwner},
			want: "_ownerid_value",
		},
		{
			name: "non-owner attribute",
			attr: models.EntityAttribute{LogicalName: "name", AttributeType: models.AttributeTypePicklist},
			want: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := queryNameForAttribute(tt.attr)
			if got != tt.want {
				t.Errorf("getQueryFieldName(%+v) = %q, want %q", tt.attr, got, tt.want)
			}
		})
	}
}

func TestGetRecordFieldName(t *testing.T) {
	const formattedSuffix = "@OData.Community.Display.V1.FormattedValue"

	tests := []struct {
		name string
		attr models.EntityAttribute
		want string
	}{
		{
			name: "owner attribute",
			attr: models.EntityAttribute{LogicalName: "ownerid", AttributeType: models.AttributeTypeOwner},
			want: "_ownerid_value" + formattedSuffix,
		},
		{
			name: "picklist attribute",
			attr: models.EntityAttribute{LogicalName: "statuscode", AttributeType: models.AttributeTypePicklist},
			want: "statuscode" + formattedSuffix,
		},
		{
			name: "state attribute",
			attr: models.EntityAttribute{LogicalName: "statecode", AttributeType: models.AttributeTypeState},
			want: "statecode" + formattedSuffix,
		},
		{
			name: "datetime attribute",
			attr: models.EntityAttribute{LogicalName: "createdon", AttributeType: models.AttributeTypeDateTime},
			want: "createdon" + formattedSuffix,
		},
		{
			name: "plain attribute",
			attr: models.EntityAttribute{LogicalName: "name"},
			want: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := recordValueNameForAttribute(tt.attr)
			if got != tt.want {
				t.Errorf("getRecordFieldName(%+v) = %q, want %q", tt.attr, got, tt.want)
			}
		})
	}
}

func TestGetSelectFields(t *testing.T) {
	attributes := []models.EntityAttribute{
		{LogicalName: "name", DisplayName: "Name"},
		{LogicalName: "ownerid", DisplayName: "Owner", AttributeType: models.AttributeTypeOwner},
	}

	tests := []struct {
		name           string
		selectedFields string
		want           []DisplayField
	}{
		{
			name:           "plain field",
			selectedFields: "name",
			want: []DisplayField{
				{DisplayName: "Name", QueryName: "name", RecordName: "name"},
			},
		},
		{
			name:           "owner field given as lookup value",
			selectedFields: "_ownerid_value",
			want: []DisplayField{
				{
					DisplayName: "Owner",
					QueryName:   "_ownerid_value",
					RecordName:  "_ownerid_value@OData.Community.Display.V1.FormattedValue",
				},
			},
		},
		{
			name:           "unknown field is skipped",
			selectedFields: "name,doesnotexist",
			want: []DisplayField{
				{DisplayName: "Name", QueryName: "name", RecordName: "name"},
			},
		},
		{
			name:           "blank entries are skipped",
			selectedFields: "name, ,",
			want: []DisplayField{
				{DisplayName: "Name", QueryName: "name", RecordName: "name"},
			},
		},
		{
			name:           "nothing selected",
			selectedFields: "",
			want:           nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSelectFields(attributes, tt.selectedFields)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSelectFields(_, %q) = %+v, want %+v", tt.selectedFields, got, tt.want)
			}
		})
	}
}

func TestGetDefaultTableFields(t *testing.T) {
	attributes := []models.EntityAttribute{
		{LogicalName: "id", IsPrimaryId: true},
		{LogicalName: "logicalid", IsPrimaryId: true, IsLogical: true},
		{LogicalName: "name", DisplayName: "Name", IsPrimaryName: true},
		{LogicalName: "createdon", AttributeType: models.AttributeTypeDateTime},
		{LogicalName: "statecode", AttributeType: models.AttributeTypeState},
		{LogicalName: "ownerid", AttributeType: models.AttributeTypeOwner},
		{LogicalName: "somethingelse"},
	}

	want := []DisplayField{
		{DisplayName: "ID", QueryName: "id", RecordName: "id"},
		{DisplayName: "Name", QueryName: "name", RecordName: "name"},
		{
			DisplayName: "Created On",
			QueryName:   "createdon",
			RecordName:  "createdon@OData.Community.Display.V1.FormattedValue",
		},
		{
			DisplayName: "State",
			QueryName:   "statecode",
			RecordName:  "statecode@OData.Community.Display.V1.FormattedValue",
		},
		{
			DisplayName: "Owner",
			QueryName:   "_ownerid_value",
			RecordName:  "_ownerid_value@OData.Community.Display.V1.FormattedValue",
		},
	}

	got := defaultDisplayFields(attributes)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("getDefaultTableFields() = %+v, want %+v", got, want)
	}
}
