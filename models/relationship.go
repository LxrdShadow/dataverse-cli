package models

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

type EntityRelationships struct {
	OneToManyRelationships  []OneToManyRelationship  `json:"OneToManyRelationships"`
	ManyToOneRelationships  []ManyToOneRelationship  `json:"ManyToOneRelationships"`
	ManyToManyRelationships []ManyToManyRelationship `json:"ManyToManyRelationships"`
}

func (r EntityRelationships) Rows(currentEntity string) []RelationshipRow {
	rows := make([]RelationshipRow, 0, len(r.OneToManyRelationships)+len(r.ManyToOneRelationships)+len(r.ManyToManyRelationships))

	for _, rel := range r.OneToManyRelationships {
		rows = append(rows, RelationshipRow{
			Type:          "OneToMany",
			SchemaName:    rel.SchemaName,
			RelatedEntity: rel.ReferencedEntity,
			FromAttribute: rel.ReferencingAttribute,
			ToAttribute:   rel.ReferencedAttribute,
		})
	}
	for _, rel := range r.ManyToOneRelationships {
		rows = append(rows, RelationshipRow{
			Type:          "ManyToOne",
			SchemaName:    rel.SchemaName,
			RelatedEntity: rel.ReferencingEntity,
			FromAttribute: rel.ReferencedAttribute,
			ToAttribute:   rel.ReferencingAttribute,
		})
	}
	for _, rel := range r.ManyToManyRelationships {
		rows = append(rows, relationshipRowFromManyToMany(rel, currentEntity))
	}
	return rows
}

func relationshipRowFromManyToMany(rel ManyToManyRelationship, currentEntity string) RelationshipRow {
	switch currentEntity {
	case rel.Entity1LogicalName:
		return RelationshipRow{
			Type:          "ManyToMany",
			SchemaName:    rel.SchemaName,
			RelatedEntity: rel.Entity2LogicalName,
			FromAttribute: rel.Entity1IntersectAttribute,
			ToAttribute:   rel.Entity2IntersectAttribute,
		}
	case rel.Entity2LogicalName:
		return RelationshipRow{
			Type:          "ManyToMany",
			SchemaName:    rel.SchemaName,
			RelatedEntity: rel.Entity1LogicalName,
			FromAttribute: rel.Entity2IntersectAttribute,
			ToAttribute:   rel.Entity1IntersectAttribute,
		}
	default:
		return RelationshipRow{
			Type:          "ManyToMany",
			SchemaName:    rel.SchemaName,
			RelatedEntity: "unknown",
			FromAttribute: "unknown",
			ToAttribute:   "unknown",
		}
	}
}
