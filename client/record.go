package client

import (
	"encoding/json"

	"dvc/models"
)

func (client *DataverseClient) ListRecords(entitySetName string) ([]models.Record, error) {
	requestUrl := client.ApiURL + "/" + entitySetName
	body, err := client.get(requestUrl)
	if err != nil {
		return nil, err
	}

	var recordResponse models.ListRecordsResponse
	if err := json.Unmarshal(body, &recordResponse); err != nil {
		return nil, err
	}

	return recordResponse.Value, nil
}
