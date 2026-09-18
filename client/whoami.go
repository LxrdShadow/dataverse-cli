package client

import (
	"dvc/models"
	"encoding/json"
)

func (client *DataverseClient) WhoAmI() (models.WhoAmIResponse, error) {
	url := client.APIURL + "/WhoAmI"

	body, err := client.get(url, nil)
	if err != nil {
		return models.WhoAmIResponse{}, err
	}

	var responseBody models.WhoAmIResponse
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return models.WhoAmIResponse{}, err
	}

	return responseBody, nil
}
