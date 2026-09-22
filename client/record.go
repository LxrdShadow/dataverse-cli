package client

import (
	"encoding/json"
	"net/url"
	"strconv"

	"dvc/models"
)

func (client *DataverseClient) ListRecords(entitySetName string, options models.ListRecordsOptions) ([]models.Record, error) {
	requestURL := client.endpoint(entitySetName)
	query := getQueryString(options.Query)
	if len(query) > 0 {
		requestURL.RawQuery = query.Encode()
	}

	headers := make(map[string]string)
	if options.Query.MaxPageSize != 0 {
		headers["Prefer"] = "odata.maxpagesize=" + strconv.Itoa(options.Query.MaxPageSize)
	}

	var records []models.Record

	for requestURL != nil && requestURL.String() != "" {
		body, err := client.get(requestURL.String(), headers)
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
		if recordResponse.NextLink != "" {
			requestURL, err = url.Parse(recordResponse.NextLink)
			if err != nil {
				return nil, err
			}
		} else {
			requestURL = nil
		}
	}

	return records, nil
}

func (client *DataverseClient) GetRecord(entitySetName string, recordID string, options models.GetRecordOptions) (models.Record, error) {
	requestURL := client.endpoint(entitySetName + "(" + recordID + ")")
	if options.Select != "" {
		query := requestURL.Query()
		query.Set("$select", options.Select)
		requestURL.RawQuery = query.Encode()
	}

	body, err := client.get(requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	var record models.Record
	if err := json.Unmarshal(body, &record); err != nil {
		return nil, err
	}
	return record, nil
}

func getQueryString(options models.RecordQueryOptions) url.Values {
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
	return query
}
