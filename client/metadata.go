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

	var response struct {
		Value []models.Entity `json:"value"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse entity list: %w", err)
	}

	return response.Value, nil
}

func (client *DataverseClient) GetTable(logicalName string, includeAttributes bool) (*models.Entity, error) {
	logicalNameFilter := fmt.Sprintf("LogicalName eq '%s'", logicalName)
	requestUrl := buildEntityListURL(client.APIURL, models.TableScopeAll, models.TableManagementAll, includeAttributes, []string{logicalNameFilter})

	body, err := client.get(requestUrl, nil)
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
	baseURL := client.APIURL + "/EntityDefinitions" + "(LogicalName='" + logicalName + "')" + "/Attributes"
	url := baseURL + "?$select=" + constants.DefaultAttributesInfo

	body, err := client.get(url, nil)
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
