package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"dvc/constants"
	"dvc/models"
)

func (client *DataverseClient) ListTables(scope models.TableScope, management models.TableManagement) ([]models.Entity, error) {
	requestUrl := buildEntityListURL(client.APIURL, scope, management, false, nil)

	body, err := client.get(requestUrl, nil)
	if err != nil {
		return nil, err
	}

	var response models.EntityDefinitionsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	var entities []models.Entity
	for _, definition := range response.Value {
		entities = append(entities, models.Entity{
			LogicalName:   definition.LogicalName,
			DisplayName:   definition.DisplayName.UserLocalizedLabel.Label,
			EntitySetName: definition.EntitySetName,
			IsCustom:      definition.IsCustomEntity,
			IsManaged:     definition.IsManaged,
		})
	}
	return entities, nil
}

func (client *DataverseClient) GetTable(logicalName string, includeAttributes bool) (*models.Entity, error) {
	logicalNameFilter := fmt.Sprintf("LogicalName eq '%s'", logicalName)
	requestUrl := buildEntityListURL(client.APIURL, models.TableScopeAll, models.TableManagementAll, includeAttributes, []string{logicalNameFilter})

	body, err := client.get(requestUrl, nil)
	if err != nil {
		return nil, err
	}

	var response models.EntityDefinitionsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response.Value) == 0 {
		message := fmt.Sprintf("entity '%s' not found", logicalName)
		if len(logicalName) > 0 && logicalName[len(logicalName)-1] == 's' {
			message = fmt.Sprintf("%s. Did you mean '%s'?", message, logicalName[:len(logicalName)-1])
		}
		return nil, fmt.Errorf("%s", message)
	}

	entity := &models.Entity{
		LogicalName:   response.Value[0].LogicalName,
		DisplayName:   response.Value[0].DisplayName.UserLocalizedLabel.Label,
		EntitySetName: response.Value[0].EntitySetName,
		IsCustom:      response.Value[0].IsCustomEntity,
		IsManaged:     response.Value[0].IsManaged,
		Attributes:    extractEntityAttributes(response.Value[0].Attributes),
	}

	return entity, nil
}

func (client *DataverseClient) ListEntityAttributes(logicalName string) ([]models.EntityAttribute, error) {
	baseURL := client.APIURL + "/EntityDefinitions" + "(LogicalName='" + logicalName + "')" + "/Attributes"
	url := baseURL + "?$select=" + constants.DefaultAttributesInfo

	body, err := client.get(url, nil)
	if err != nil {
		return nil, err
	}

	var response models.EntityAttributesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	attributes := extractEntityAttributes(response.Value)
	return attributes, nil
}

func buildEntityListURL(baseURL string, scope models.TableScope, management models.TableManagement, includeAttributes bool, customFilters []string) string {
	query := url.Values{}
	query.Set("$select", constants.DefaultEntityAttributes)

	if includeAttributes {
		query.Set("$expand", "Attributes($select="+constants.DefaultAttributesInfo+")")
	}

	filters := append([]string{}, customFilters...)

	switch scope {
	case models.TableScopeSystem:
		filters = append(filters, "IsCustomEntity eq false")
	case models.TableScopeCustom:
		filters = append(filters, "IsCustomEntity eq true")
	}

	switch management {
	case models.TableManagementManaged:
		filters = append(filters, "IsManaged eq true")
	case models.TableManagementUnmanaged:
		filters = append(filters, "IsManaged eq false")
	}

	if len(filters) > 0 {
		query.Set("$filter", strings.Join(filters, " and "))
	}

	return baseURL + "/EntityDefinitions?" + query.Encode()
}

func extractEntityAttributes(rawAttributes []models.RawEntityAttribute) []models.EntityAttribute {
	var attributes []models.EntityAttribute
	for _, attr := range rawAttributes {
		attributes = append(attributes, models.EntityAttribute{
			LogicalName:   attr.LogicalName,
			DisplayName:   attr.DisplayName.UserLocalizedLabel.Label,
			IsPrimaryName: attr.IsPrimaryName,
			IsPrimaryId:   attr.IsPrimaryId,
			IsLogical:     attr.IsLogical,
			AttributeType: attr.AttributeType,
		})
	}
	return attributes
}
