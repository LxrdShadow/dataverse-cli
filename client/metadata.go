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
	requestURL := buildEntityListURL(client.endpoint(), scope, management, false, nil)

	body, err := client.get(requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Value []models.Entity `json:"value"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse entity list: %w", err)
	}

	return response.Value, nil
}

func (client *DataverseClient) GetTable(logicalName string, includeAttributes bool) (*models.Entity, error) {
	logicalNameFilter := fmt.Sprintf("LogicalName eq %s", odataStringLiteral(logicalName))
	requestURL := buildEntityListURL(client.endpoint(), models.TableScopeAll, models.TableManagementAll, includeAttributes, []string{logicalNameFilter})

	body, err := client.get(requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Value []models.Entity `json:"value"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse entity: %w", err)
	}

	if len(response.Value) == 0 {
		return nil, fmt.Errorf("entity '%s' not found", logicalName)
	}

	return &response.Value[0], nil
}

func (client *DataverseClient) ListEntityAttributes(logicalName string) ([]models.EntityAttribute, error) {
	requestURL := client.entityAttributesURL(logicalName)
	query := requestURL.Query()
	query.Set("$select", constants.DefaultAttributesInfo)
	requestURL.RawQuery = query.Encode()

	body, err := client.get(requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Value []models.EntityAttribute `json:"Attributes"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse entity attributes: %w", err)
	}

	return response.Value, nil
}

func (client *DataverseClient) entityAttributesURL(logicalName string) *url.URL {
	endpoint := *client.apiURL
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") +
		"/" + entityDefinitionPath(logicalName)

	query := endpoint.Query()
	query.Set("$select", constants.DefaultAttributesInfo)
	endpoint.RawQuery = query.Encode()

	return &endpoint
}

func buildEntityListURL(baseURL *url.URL, scope models.TableScope, management models.TableManagement, includeAttributes bool, customFilters []string) *url.URL {
	query := baseURL.Query()
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

	requestURL := baseURL.JoinPath("EntityDefinitions")
	requestURL.RawQuery = query.Encode()
	return requestURL
}
