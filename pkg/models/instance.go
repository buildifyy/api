package models

type Instance struct {
	BasicInformations []InstanceBasicInformation `bson:"basicInformation" json:"basicInformation"`
	Attributes        []InstanceAttribute        `bson:"attributes" json:"attributes"`
	MetricTypes       []InstanceMetricType       `bson:"metricTypes" json:"metricTypes"`
	TenantID          string                     `bson:"tenantId" json:"tenantId"`
}

type InstanceFormMetaData struct {
	BasicInformation InstanceMetaData `json:"basicInformation"`
	Attributes       InstanceMetaData `json:"attributes"`
}

type InstanceMetaData struct {
	Fields []InstanceMetaDataFields `json:"fields"`
}

type InstanceMetaDataFields struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	InfoText   string `json:"infoText"`
	TypeLabel  string `json:"typeLabel"`
	Type       string `json:"type"`
	IsRequired bool   `json:"isRequired"`
	IsHidden   bool   `json:"isHidden"`
}

type InstanceBasicInformation struct {
	ID    string `bson:"id" json:"id"`
	Name  string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
}

type InstanceAttribute struct {
	ID         string      `bson:"id" json:"id"`
	Name       string      `bson:"name" json:"name"`
	DataType   string      `bson:"dataType" json:"dataType"`
	IsRequired bool        `bson:"isRequired" json:"isRequired"`
	IsHidden   bool        `bson:"isHidden" json:"isHidden"`
	Value      interface{} `bson:"value" json:"value"`
}

type InstanceMetricType struct {
	ID      string           `bson:"id" json:"id"`
	Name    string           `bson:"name" json:"name"`
	Metrics []InstanceMetric `bson:"metrics" json:"metrics"`
}

type InstanceMetric struct {
	ID     string      `bson:"id" json:"id"`
	Name   string      `bson:"name" json:"name"`
	Manual interface{} `bson:"manual" json:"manual"`
}
