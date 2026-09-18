package client

import (
	"dvc/errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DataverseClient struct {
	BaseURL        string
	APIURL         string
	Token          string
	RequestTimeout time.Duration
	HTTP           *http.Client
}

func NewDataverseClient(baseURL string, token string, timeout time.Duration) *DataverseClient {
	// TODO: Setup timeout on HTTP client
	return &DataverseClient{
		BaseURL:        baseURL,
		APIURL:         strings.TrimRight(baseURL, "/") + "/api/data/v9.2",
		Token:          token,
		RequestTimeout: timeout,
		HTTP:           &http.Client{Timeout: timeout},
	}
}

func (client *DataverseClient) get(url string, headers map[string]string) ([]byte, error) {
	return client.request(http.MethodGet, url, headers)
}

func (client *DataverseClient) request(method string, url string, headers map[string]string) ([]byte, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+client.Token)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OData-Version", "4.0")
	request.Header.Set("OData-MaxVersion", "4.0")
	request.Header.Set("Prefer", "odata.include-annotations=*")

	for key, value := range headers {
		request.Header.Set(key, value)
	}

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
		if response.StatusCode == 401 {
			return nil, errors.ErrUnauthorized
		}

		var errorResponse errors.ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err == nil {
			return body, errorResponse.ErrorValue
		}
		return body, fmt.Errorf("Unexpected status code: %d", response.StatusCode)
	}

	return body, nil
}
