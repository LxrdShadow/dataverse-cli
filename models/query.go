package models

type QueryOptions struct {
	Select      string
	Filter      string
	OrderBy     string
	Expand      string
	MaxPageSize int
	Top         int
}

type ListOptions struct {
	Output OutputFormat
	Query  QueryOptions
}

type GetOptions struct {
	Select string
	Output OutputFormat
}
