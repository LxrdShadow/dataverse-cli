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
	Type             string
	SchemaName       string
	CurrentEntity    string
	CurrentAttribute string
	RelatedEntity    string
	RelatedAttribute string
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
			Type:             "OneToMany",
			SchemaName:       rel.SchemaName,
			CurrentEntity:    rel.ReferencedEntity,
			CurrentAttribute: rel.ReferencedAttribute,
			RelatedEntity:    rel.ReferencingEntity,
			RelatedAttribute: rel.ReferencingAttribute,
		})
	}
	for _, rel := range r.ManyToOneRelationships {
		rows = append(rows, RelationshipRow{
			Type:             "ManyToOne",
			SchemaName:       rel.SchemaName,
			CurrentEntity:    rel.ReferencingEntity,
			CurrentAttribute: rel.ReferencingAttribute,
			RelatedEntity:    rel.ReferencedEntity,
			RelatedAttribute: rel.ReferencedAttribute,
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
			Type:             "ManyToMany",
			SchemaName:       rel.SchemaName,
			CurrentEntity:    rel.Entity1LogicalName,
			CurrentAttribute: rel.Entity1IntersectAttribute,
			RelatedEntity:    rel.Entity2LogicalName,
			RelatedAttribute: rel.Entity2IntersectAttribute,
		}
	case rel.Entity2LogicalName:
		return RelationshipRow{
			Type:             "ManyToMany",
			SchemaName:       rel.SchemaName,
			CurrentEntity:    rel.Entity2LogicalName,
			CurrentAttribute: rel.Entity2IntersectAttribute,
			RelatedEntity:    rel.Entity1LogicalName,
			RelatedAttribute: rel.Entity1IntersectAttribute,
		}
	default:
		return RelationshipRow{
			Type:             "ManyToMany",
			SchemaName:       rel.SchemaName,
			CurrentEntity:    "unknown",
			CurrentAttribute: "unknown",
			RelatedEntity:    "unknown",
			RelatedAttribute: "unknown",
		}
	}
}
