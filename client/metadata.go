package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"dvc/constants"
	"dvc/models"
)

func (c *DataverseClient) ListEntities(scope models.EntityScope, management models.EntityManagement) ([]models.Entity, error) {
	requestURL := buildEntityListURL(c.endpoint(), scope, management, false, nil)

	body, err := c.getURL(requestURL.String(), nil)
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

func (c *DataverseClient) GetEntity(logicalName string, includeAttributes bool) (*models.Entity, error) {
	logicalNameFilter := fmt.Sprintf("LogicalName eq %s", odataStringLiteral(logicalName))
	requestURL := buildEntityListURL(c.endpoint(), models.EntityScopeAll, models.EntityManagementAll, includeAttributes, []string{logicalNameFilter})

	body, err := c.getURL(requestURL.String(), nil)
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

func (c *DataverseClient) ListEntityAttributes(logicalName string) ([]models.EntityAttribute, error) {
	requestURL := c.entityAttributesURL(logicalName)
	query := requestURL.Query()
	query.Set("$select", constants.DefaultAttributesInfo)
	requestURL.RawQuery = query.Encode()

	body, err := c.getURL(requestURL.String(), nil)
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

func (c *DataverseClient) entityAttributesURL(logicalName string) *url.URL {
	endpoint := *c.apiURL
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") +
		"/" + entityDefinitionPath(logicalName)

	query := endpoint.Query()
	query.Set("$select", constants.DefaultAttributesInfo)
	endpoint.RawQuery = query.Encode()

	return &endpoint
}

func buildEntityListURL(baseURL *url.URL, scope models.EntityScope, management models.EntityManagement, includeAttributes bool, customFilters []string) *url.URL {
	query := baseURL.Query()
	query.Set("$select", constants.DefaultEntityAttributes)

	if includeAttributes {
		query.Set("$expand", "Attributes($select="+constants.DefaultAttributesInfo+")")
	}

	filters := append([]string{}, customFilters...)

	switch scope {
	case models.EntityScopeSystem:
		filters = append(filters, "IsCustomEntity eq false")
	case models.EntityScopeCustom:
		filters = append(filters, "IsCustomEntity eq true")
	}

	switch management {
	case models.EntityManagementManaged:
		filters = append(filters, "IsManaged eq true")
	case models.EntityManagementUnmanaged:
		filters = append(filters, "IsManaged eq false")
	}

	if len(filters) > 0 {
		query.Set("$filter", strings.Join(filters, " and "))
	}

	requestURL := baseURL.JoinPath("EntityDefinitions")
	requestURL.RawQuery = query.Encode()
	return requestURL
}
