package models

type WhoAmIResponse struct {
	BusinessUnitId string `json:"businessUnitId"`
	UserId         string `json:"userId"`
	OrganizationId string `json:"organizationId"`
}
