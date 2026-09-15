package client

import (
	"dvc/models"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func (client *DataverseClient) ListTables(scope models.TableScope, management models.TableManagement) ([]models.Entity, error) {
	requestUrl := buildURL(client.ApiURL, scope, management, nil)

	body, err := client.get(requestUrl)
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

func (client *DataverseClient) GetTable(logicalName string) (*models.Entity, error) {
	logicalNameFilter := fmt.Sprintf("LogicalName eq '%s'", logicalName)
	requestUrl := buildURL(client.ApiURL, models.TableScopeAll, models.TableManagementAll, []string{logicalNameFilter})

	body, err := client.get(requestUrl)
	if err != nil {
		return nil, err
	}

	var response models.EntityDefinitionsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response.Value) == 0 {
		return nil, fmt.Errorf("entity %s not found", logicalName)
	}

	entity := &models.Entity{
		LogicalName:   response.Value[0].LogicalName,
		DisplayName:   response.Value[0].DisplayName.UserLocalizedLabel.Label,
		EntitySetName: response.Value[0].EntitySetName,
		IsCustom:      response.Value[0].IsCustomEntity,
		IsManaged:     response.Value[0].IsManaged,
	}

	return entity, nil
}

func buildURL(baseURL string, scope models.TableScope, management models.TableManagement, customFilters []string) string {
	var builder strings.Builder
	builder.WriteString(baseURL)
	builder.WriteString("/EntityDefinitions")
	builder.WriteString("?$select=LogicalName,EntitySetName,DisplayName,IsCustomEntity,IsManaged")

	var filters []string

	if len(customFilters) > 0 {
		filters = append(filters, customFilters...)
	}

	switch scope {
	case models.TableScopeSystem:
		filters = append(filters, url.PathEscape("IsCustomEntity eq false"))
	case models.TableScopeCustom:
		filters = append(filters, url.PathEscape("IsCustomEntity eq true"))
	}

	switch management {
	case models.TableManagementManaged:
		filters = append(filters, url.PathEscape("IsManaged eq true"))
	case models.TableManagementUnmanaged:
		filters = append(filters, url.PathEscape("IsManaged eq false"))
	}

	if len(filters) > 0 {
		builder.WriteString("&$filter=")
		builder.WriteString(strings.Join(filters, url.PathEscape(" and ")))
	}

	return builder.String()
}
