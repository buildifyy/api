package models

type Template struct {
	TenantID         string           `bson:"tenantId" json:"-"`
	BasicInformation BasicInformation `bson:"basicInformation" json:"basicInformation"`
	Attributes       []Attribute      `bson:"attributes" json:"attributes"`
	MetricTypes      []MetricType     `bson:"metricTypes" json:"metricTypes"`
}

type BasicInformation struct {
	Name       string `bson:"name" json:"name"`
	Parent     string `bson:"parent" json:"parent"`
	ExternalID string `bson:"externalId" json:"externalId"`
	IsCustom   bool   `bson:"isCustom" json:"isCustom"`
}

type Attribute struct {
	Name       string `bson:"name" json:"name"`
	DataType   string `bson:"dataType" json:"dataType"`
	IsRequired bool   `bson:"isRequired" json:"isRequired"`
	IsHidden   bool   `bson:"isHidden" json:"isHidden"`
}

type MetricType struct {
	Name       string   `bson:"name" json:"name"`
	MetricType string   `bson:"metricType" json:"metricType"`
	Metrics    []Metric `bson:"metrics" json:"metrics"`
}

type Metric struct {
	Name         string      `bson:"name" json:"name"`
	IsManual     bool        `bson:"isManual" json:"isManual"`
	Value        interface{} `bson:"value" json:"value"`
	IsCalculated bool        `bson:"isCalculated" json:"isCalculated"`
	IsSourced    bool        `bson:"isSourced" json:"isSourced"`
}
