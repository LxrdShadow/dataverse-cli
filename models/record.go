package models

type Record map[string]any

type ListRecordsResponse struct {
	Value []Record `json:"value"`
}

type QueryOptions struct {
	Select  string
	Filter  string
	OrderBy string
	Expand  string
	Top     int
}
