package client

import (
	"dvc/errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type DataverseClient struct {
	BaseURL string
	ApiURL  string
	Token   string
	HTTP    *http.Client
}

func NewDataverseClient(baseURL string, token string) *DataverseClient {
	return &DataverseClient{
		BaseURL: baseURL,
		ApiURL:  strings.TrimRight(baseURL, "/") + "/api/data/v9.2",
		Token:   token,
		HTTP:    &http.Client{},
	}
}

func (client *DataverseClient) get(url string) ([]byte, error) {
	return client.request(http.MethodGet, url)
}

func (client *DataverseClient) request(method string, url string) ([]byte, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+client.Token)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OData-Version", "4.0")
	request.Header.Set("OData-MaxVersion", "4.0")

	response, err := client.HTTP.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var errorResponse errors.ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err == nil {
			return nil, fmt.Errorf("%s", errorResponse.ErrorValue.Message)
		}
		return nil, errorResponse
	}

	return body, nil
}
