package models

type EntityScope string
type EntityManagement string

const (
	EntityScopeAll    EntityScope = "all"
	EntityScopeCustom EntityScope = "custom"
	EntityScopeSystem EntityScope = "system"
)

const (
	EntityManagementAll       EntityManagement = "all"
	EntityManagementManaged   EntityManagement = "managed"
	EntityManagementUnmanaged EntityManagement = "unmanaged"
)
