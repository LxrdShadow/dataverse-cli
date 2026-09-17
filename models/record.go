package models

type Record map[string]any

type ListRecordsResponse struct {
	NextLink string   `json:"@odata.nextLink"`
	Value    []Record `json:"value"`
}
