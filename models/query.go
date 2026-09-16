package models

type QueryOptions struct {
	Select  string
	Filter  string
	OrderBy string
	Expand  string
	Top     int
}
