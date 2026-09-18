package client

import (
	"encoding/json"
	"net/url"
	"strconv"

	"dvc/models"
)

func (client *DataverseClient) ListRecords(entitySetName string, options models.ListOptions) ([]models.Record, error) {
	requestUrl := client.APIURL + "/" + entitySetName
	queryString := getQueryString(options.Query)
	if queryString != "" {
		requestUrl += "?" + queryString
	}

	headers := make(map[string]string)
	if options.Query.MaxPageSize != 0 {
		headers["Prefer"] = "odata.maxpagesize=" + strconv.Itoa(options.Query.MaxPageSize)
	}

	var records []models.Record

	for requestUrl != "" {
		body, err := client.get(requestUrl, headers)
		if err != nil {
			return nil, err
		}

		var recordResponse models.ListRecordsResponse
		if err := json.Unmarshal(body, &recordResponse); err != nil {
			return nil, err
		}

		if options.Query.Top != 0 {
			remaining := options.Query.Top - len(records)
			if len(recordResponse.Value) > remaining {
				records = append(records, recordResponse.Value[:remaining]...)
				break
			}
		}

		records = append(records, recordResponse.Value...)
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
