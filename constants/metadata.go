package constants

import "time"

const (
	DefaultEntityAttributes = "LogicalName,EntitySetName,DisplayName,IsCustomEntity,IsManaged,PrimaryIdAttribute,PrimaryNameAttribute"
	DefaultAttributesInfo   = "LogicalName,DisplayName,IsPrimaryName,IsPrimaryId,IsLogical,RequiredLevel,AttributeType"

	DefaultTimeout = 10 * time.Second
)
