package models

type Template struct {
	TenantID         string                   `bson:"tenantId" json:"-"`
	BasicInformation TemplateBasicInformation `bson:"basicInformation" json:"basicInformation"`
	Attributes       []TemplateAttribute      `bson:"attributes" json:"attributes"`
	MetricTypes      []TemplateMetricType     `bson:"metricTypes" json:"metricTypes"`
}

type TemplateBasicInformation struct {
	Name       string `bson:"name" json:"name"`
	Parent     string `bson:"parent" json:"parent"`
	ExternalID string `bson:"externalId" json:"externalId"`
	IsCustom   bool   `bson:"isCustom" json:"isCustom"`
}

type TemplateAttribute struct {
	ID             string `bson:"id" json:"id"`
	Name           string `bson:"name" json:"name"`
	DataType       string `bson:"dataType" json:"dataType"`
	IsRequired     bool   `bson:"isRequired" json:"isRequired"`
	IsHidden       bool   `bson:"isHidden" json:"isHidden"`
	OwningTemplate string `bson:"owningTemplate" json:"owningTemplate"`
}

type TemplateMetricType struct {
	ID             string           `bson:"id" json:"id"`
	Name           string           `bson:"name" json:"name"`
	MetricType     string           `bson:"metricType" json:"metricType"`
	Metrics        []TemplateMetric `bson:"metrics" json:"metrics"`
	OwningTemplate string           `bson:"owningTemplate" json:"owningTemplate"`
}

type TemplateMetric struct {
	ID           string      `bson:"id" json:"id"`
	Name         string      `bson:"name" json:"name"`
	IsManual     bool        `bson:"isManual" json:"isManual"`
	Value        interface{} `bson:"value" json:"value"`
	IsCalculated bool        `bson:"isCalculated" json:"isCalculated"`
	IsSourced    bool        `bson:"isSourced" json:"isSourced"`
}
