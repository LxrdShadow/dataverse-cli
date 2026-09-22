package constants

import "time"

const (
	DefaultEntityAttributes = "LogicalName,EntitySetName,DisplayName,IsCustomEntity,IsManaged,PrimaryIdAttribute,PrimaryNameAttribute,Attributes"
	DefaultAttributesInfo   = "LogicalName,DisplayName,IsPrimaryName,IsPrimaryId,IsLogical,AttributeType"

	DefaultTimeout = 10 * time.Second
)
