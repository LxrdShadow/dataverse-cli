package models

type QueryOptions struct {
	Select      string
	Filter      string
	OrderBy     string
	Expand      string
	MaxPageSize int
	Top         int
}
