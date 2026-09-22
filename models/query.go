package models

type RecordQueryOptions struct {
	Select      string
	Filter      string
	OrderBy     string
	Expand      string
	MaxPageSize int
	Top         int
}

type ListRecordsOptions struct {
	Output OutputFormat
	Query  RecordQueryOptions
}

type GetRecordOptions struct {
	Select       string
	Output       OutputFormat
	DisplayEmpty bool
}
