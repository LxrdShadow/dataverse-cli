package models

type EntityDefinitionsResponse struct {
	Value []EntityDefinition `json:"value"`
}

type EntityDefinition struct {
	LogicalName    string               `json:"LogicalName"`
	EntitySetName  string               `json:"EntitySetName"`
	DisplayName    DisplayName          `json:"DisplayName"`
	IsCustomEntity bool                 `json:"IsCustomEntity"`
	IsManaged      bool                 `json:"IsManaged"`
	Attributes     []RawEntityAttribute `json:"Attributes"`
}

type DisplayName struct {
	UserLocalizedLabel LocalizedLabel `json:"UserLocalizedLabel"`
}

type LocalizedLabel struct {
	Label string `json:"Label"`
}

type Entity struct {
	LogicalName   string
	DisplayName   string
	EntitySetName string
	IsCustom      bool
	IsManaged     bool
	Attributes    []EntityAttribute
}

type TableScope string

const (
	TableScopeAll    TableScope = "all"
	TableScopeCustom TableScope = "custom"
	TableScopeSystem TableScope = "system"
)

type TableManagement string

const (
	TableManagementAll       TableManagement = "all"
	TableManagementManaged   TableManagement = "managed"
	TableManagementUnmanaged TableManagement = "unmanaged"
)

type EntityAttributesResponse struct {
	Value []RawEntityAttribute `json:"value"`
}

type RawEntityAttribute struct {
	LogicalName   string      `json:"LogicalName"`
	DisplayName   DisplayName `json:"DisplayName"`
	IsPrimaryName bool        `json:"IsPrimaryName"`
	IsPrimaryId   bool        `json:"IsPrimaryId"`
	IsLogical     bool        `json:"IsLogical"`
}

type EntityAttribute struct {
	LogicalName   string
	DisplayName   string
	IsPrimaryName bool
	IsPrimaryId   bool
	IsLogical     bool
}

type GenericTableStructure struct {
	Id          string
	PrimaryName string
	CreatedOn   string
	State       string
	Owner       string
}
