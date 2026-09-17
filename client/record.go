package client

import (
	"encoding/json"
	"net/url"
	"strconv"

	"dvc/models"
)

func (client *DataverseClient) ListRecords(entitySetName string, options models.QueryOptions) ([]models.Record, error) {
	requestUrl := client.ApiURL + "/" + entitySetName
	queryString := getQueryString(options)
	if queryString != "" {
		requestUrl += "?" + queryString
	}

	var records []models.Record
	next := true

	for next {
		body, err := client.get(requestUrl)
		if err != nil {
			return nil, err
		}

		var recordResponse models.ListRecordsResponse
		if err := json.Unmarshal(body, &recordResponse); err != nil {
			return nil, err
		}

		records = append(records, recordResponse.Value...)
		next = recordResponse.NextLink != ""
		requestUrl = recordResponse.NextLink
	}

	return records, nil
}

func getQueryString(options models.QueryOptions) string {
	query := url.Values{}
	if len(options.Select) > 0 {
		query.Set("$select", options.Select)
	}
	if len(options.Filter) > 0 {
		query.Set("$filter", options.Filter)
	}
	if options.OrderBy != "" {
		query.Set("$orderby", options.OrderBy)
	}
	if options.Expand != "" {
		query.Set("$expand", options.Expand)
	}
	if options.Top != 0 {
		query.Set("$top", strconv.Itoa(options.Top))
	}
	return query.Encode()
}
