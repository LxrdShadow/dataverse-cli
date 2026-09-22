package models

type WhoAmIResponse struct {
	BusinessUnitID string `json:"businessUnitId"`
	UserID         string `json:"userId"`
	OrganizationID string `json:"organizationId"`
}
