package client

import (
	"dvc/models"
	"encoding/json"
	"fmt"
)

func (client *DataverseClient) WhoAmI() (models.WhoAmIResponse, error) {
	url := client.APIURL + "/WhoAmI"

	body, err := client.get(url, nil)
	if err != nil {
		return models.WhoAmIResponse{}, err
	}

	var responseBody models.WhoAmIResponse
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return models.WhoAmIResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}

	return responseBody, nil
}
