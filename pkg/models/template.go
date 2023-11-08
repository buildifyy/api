package models

type Template struct {
	TenantID         string                   `bson:"tenantId" json:"-"`
	BasicInformation TemplateBasicInformation `bson:"basicInformation" json:"basicInformation"`
	Attributes       []TemplateAttribute      `bson:"attributes" json:"attributes"`
	Metrics          []TemplateMetric         `bson:"metrics" json:"metrics"`
}

type TemplateBasicInformation struct {
	Name         string `bson:"name" json:"name"`
	Parent       string `bson:"parent" json:"parent"`
	ExternalID   string `bson:"externalId" json:"externalId"`
	IsCustom     bool   `bson:"isCustom" json:"isCustom"`
	RootTemplate string `bson:"rootTemplate" json:"rootTemplate"`
}

type TemplateAttribute struct {
	ID             string `bson:"id" json:"id"`
	Name           string `bson:"name" json:"name"`
	DataType       string `bson:"dataType" json:"dataType"`
	IsRequired     bool   `bson:"isRequired" json:"isRequired"`
	IsHidden       bool   `bson:"isHidden" json:"isHidden"`
	OwningTemplate string `bson:"owningTemplate" json:"owningTemplate"`
}

type TemplateMetric struct {
	ID             string      `bson:"id" json:"id"`
	Name           string      `bson:"name" json:"name"`
	MetricType     string      `bson:"metricType" json:"metricType"`
	Unit           string      `bson:"unit" json:"unit"`
	IsManual       bool        `bson:"isManual" json:"isManual"`
	Value          interface{} `bson:"value" json:"value"`
	IsCalculated   bool        `bson:"isCalculated" json:"isCalculated"`
	IsSourced      bool        `bson:"isSourced" json:"isSourced"`
	OwningTemplate string      `bson:"owningTemplate" json:"owningTemplate"`
}
