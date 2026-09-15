package models

type Record map[string]any

type ListRecordsResponse struct {
	Value []Record `json:"value"`
}

type QueryOptions struct {
	Select  []string
	Filters []string
	OrderBy string
	Expand  string
	Top     int
}
