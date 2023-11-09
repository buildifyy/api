package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Instance struct {
	BasicInformation InstanceBasicInformation `bson:"basicInformation" json:"basicInformation"`
	Attributes       []InstanceAttribute      `bson:"attributes" json:"attributes"`
	Metrics          []InstanceMetric         `bson:"metrics" json:"metrics"`
	Relationships    []InstanceRelationship   `bson:"relationships" json:"relationships"`
	TenantID         string                   `bson:"tenantId" json:"tenantId"`
}

type InstanceBasicInformation struct {
	Parent       string `bson:"parent" json:"parent"`
	ExternalId   string `bson:"externalId" json:"externalId"`
	Name         string `bson:"name" json:"name"`
	IsCustom     bool   `bson:"isCustom" json:"isCustom"`
	RootTemplate string `bson:"rootTemplate" json:"rootTemplate"`
}

type InstanceAttribute struct {
	ID    string      `bson:"id" json:"id"`
	Value interface{} `bson:"value" json:"value"`
}

type InstanceMetric struct {
	ID              string      `bson:"id" json:"id"`
	MetricBehaviour string      `bson:"metricBehaviour" json:"metricBehaviour"`
	Value           interface{} `bson:"value" json:"value"`
}

type InstanceRelationship struct {
	ID                     string             `bson:"id" json:"id"`
	Target                 interface{}        `bson:"target" json:"target"`
	RelationshipTemplateId primitive.ObjectID `bson:"relationshipTemplateId" json:"relationshipTemplateId"`
}

type InstanceFormMetaData struct {
	BasicInformation InstanceMetaData `json:"basicInformation"`
	Attributes       InstanceMetaData `json:"attributes"`
	Metrics          InstanceMetaData `json:"metrics"`
}

type InstanceMetaData struct {
	Fields []InstanceMetaDataFields `json:"fields"`
}

type InstanceMetaDataFields struct {
	ID             string      `json:"id"`
	Label          string      `json:"label"`
	InfoText       string      `json:"infoText"`
	TypeLabel      string      `json:"typeLabel"`
	Type           string      `json:"type"`
	IsRequired     bool        `json:"isRequired"`
	IsHidden       bool        `json:"isHidden"`
	DropdownValues []string    `json:"dropdownValues"`
	ManualValue    interface{} `json:"manualValue"`
	Unit           string      `json:"unit"`
}
