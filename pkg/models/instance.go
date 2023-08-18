package models

type Instance struct {
	BasicInformation InstanceBasicInformation `bson:"basicInformation" json:"basicInformation"`
	Attributes       []InstanceAttribute      `bson:"attributes" json:"attributes"`
	MetricTypes      []InstanceMetricType     `bson:"metricTypes" json:"metricTypes"`
	TenantID         string                   `bson:"tenantId" json:"tenantId"`
}

type InstanceBasicInformation struct {
	Name       string `bson:"name" json:"name"`
	ExternalID string `bson:"externalId" json:"externalId"`
	Parent     string `bson:"parent" json:"parent"`
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
