package models

type EntityRelationships struct {
	OneToManyRelationships  []OneToManyRelationship  `json:"OneToManyRelationships"`
	ManyToOneRelationships  []ManyToOneRelationship  `json:"ManyToOneRelationships"`
	ManyToManyRelationships []ManyToManyRelationship `json:"ManyToManyRelationships"`
}

type OneToManyRelationship struct {
	SchemaName           string `json:"SchemaName"`
	ReferencedEntity     string `json:"ReferencedEntity"`
	ReferencedAttribute  string `json:"ReferencedAttribute"`
	ReferencingEntity    string `json:"ReferencingEntity"`
	ReferencingAttribute string `json:"ReferencingAttribute"`
}

type ManyToOneRelationship struct {
	SchemaName           string `json:"SchemaName"`
	ReferencedEntity     string `json:"ReferencedEntity"`
	ReferencedAttribute  string `json:"ReferencedAttribute"`
	ReferencingEntity    string `json:"ReferencingEntity"`
	ReferencingAttribute string `json:"ReferencingAttribute"`
}

type ManyToManyRelationship struct {
	SchemaName                string `json:"SchemaName"`
	Entity1LogicalName        string `json:"Entity1LogicalName"`
	Entity1IntersectAttribute string `json:"Entity1IntersectAttribute"`
	Entity2LogicalName        string `json:"Entity2LogicalName"`
	Entity2IntersectAttribute string `json:"Entity2IntersectAttribute"`
}

type RelationshipRow struct {
	Type          string
	SchemaName    string
	RelatedEntity string
	FromAttribute string
	ToAttribute   string
}
