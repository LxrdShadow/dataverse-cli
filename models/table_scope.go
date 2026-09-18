package models

type TableScope string
type TableManagement string

const (
	TableScopeAll    TableScope = "all"
	TableScopeCustom TableScope = "custom"
	TableScopeSystem TableScope = "system"
)

const (
	TableManagementAll       TableManagement = "all"
	TableManagementManaged   TableManagement = "managed"
	TableManagementUnmanaged TableManagement = "unmanaged"
)
