package models

type QueryOptions struct {
	Select  string
	Filter  string
	OrderBy string
	Expand  string
	Top     int
}

type ListOptions struct {
	Output OutputFormat
	Query  QueryOptions
}
