package constants

import "time"

const (
	DefaultEntityAttributes  = "LogicalName,EntitySetName,DisplayName,IsCustomEntity,IsManaged,PrimaryIdAttribute,PrimaryNameAttribute"
	DefaultAttributesInfo    = "LogicalName,DisplayName,IsPrimaryName,IsPrimaryId,IsLogical,RequiredLevel,AttributeType"
	DefaultManyToManyInfo    = "RelationshipType,SchemaName,Entity1LogicalName,Entity1IntersectAttribute,Entity2LogicalName,Entity2IntersectAttribute"
	DefaultOneToManyInfo     = "RelationshipType,SchemaName,ReferencedEntity,ReferencedAttribute,ReferencingEntity,ReferencingAttribute"
	DefaultRelationshipsInfo = "ManyToManyRelationships($select=" + DefaultManyToManyInfo + "),ManyToOneRelationships($select=" + DefaultOneToManyInfo + "),OneToManyRelationships($select=" + DefaultOneToManyInfo + ")"

	DefaultTimeout = 10 * time.Second
)
