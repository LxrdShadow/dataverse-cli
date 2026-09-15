package client

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"dvc/models"
)

func (client *DataverseClient) ListRecords(entitySetName string, options models.QueryOptions) ([]models.Record, error) {
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

func getQueryString(options models.QueryOptions) string {
	query := url.Values{}
	if len(options.Select) > 0 {
		query.Add("$select", strings.Join(options.Select, ","))
	}
	if len(options.Filters) > 0 {
		query.Add("$filter", strings.Join(options.Filters, " and "))
	}
	if options.OrderBy != "" {
		query.Add("$orderby=", options.OrderBy)
	}
	if options.Expand != "" {
		query.Add("$expand=", options.Expand)
	}
	if options.Top != 0 {
		query.Add("$top", strconv.Itoa(options.Top))
	}
	return query.Encode()
}
