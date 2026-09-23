package constants

import "time"

const (
	DefaultEntityAttributesSelect = "LogicalName,EntitySetName,DisplayName,IsCustomEntity,IsManaged,PrimaryIdAttribute,PrimaryNameAttribute"
	DefaultAttributesSelect       = "LogicalName,DisplayName,IsPrimaryName,IsPrimaryId,IsLogical,RequiredLevel,AttributeType"
	DefaultManyToManySelect       = "RelationshipType,SchemaName,Entity1LogicalName,Entity1IntersectAttribute,Entity2LogicalName,Entity2IntersectAttribute"
	DefaultOneToManySelect        = "RelationshipType,SchemaName,ReferencedEntity,ReferencedAttribute,ReferencingEntity,ReferencingAttribute"
	DefaultRelationshipsSelect    = "ManyToManyRelationships($select=" + DefaultManyToManySelect + "),ManyToOneRelationships($select=" + DefaultOneToManySelect + "),OneToManyRelationships($select=" + DefaultOneToManySelect + ")"

	DefaultTimeout = 10 * time.Second
)
