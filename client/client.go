package client

import (
	"dvc/errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type DataverseClient struct {
	apiURL         *url.URL
	Token          string
	RequestTimeout time.Duration
	HTTP           *http.Client
}

func NewDataverseClient(baseURL string, token string, timeout time.Duration) (*DataverseClient, error) {
	parsedURL, err := url.Parse(strings.TrimRight(baseURL, "/") + "/api/data/v9.2")
	if err != nil {
		return nil, fmt.Errorf("invalid Dataverse URL: %w", err)
	}

	return &DataverseClient{
		apiURL:         parsedURL,
		Token:          token,
		RequestTimeout: timeout,
		HTTP:           &http.Client{Timeout: timeout},
	}, nil
}

func (c *DataverseClient) endpoint(parts ...string) *url.URL {
	endpoint := *c.apiURL
	endpoint.Path = path.Join(endpoint.Path, path.Join(parts...))
	return &endpoint
}

func (c *DataverseClient) getURL(rawURL string, headers map[string]string) ([]byte, error) {
	return c.doRequest(http.MethodGet, rawURL, headers)
}

func (c *DataverseClient) doRequest(method string, rawURL string, headers map[string]string) ([]byte, error) {
	request, err := http.NewRequest(method, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+c.Token)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OData-Version", "4.0")
	request.Header.Set("OData-MaxVersion", "4.0")
	request.Header.Set("Prefer", "odata.include-annotations=*")

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := c.HTTP.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if response.StatusCode == 401 {
			return nil, errors.ErrUnauthorized
		}

		var errorResponse errors.ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err == nil {
			fmt.Println(rawURL)
			return body, errorResponse.ErrorValue
		}
		return body, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return body, nil
}
