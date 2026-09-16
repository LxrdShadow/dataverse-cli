package models

type Record map[string]any

type ListRecordsResponse struct {
	Value []Record `json:"value"`
}
