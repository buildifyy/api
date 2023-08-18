package models

type Instance struct {
	BasicInformation InstanceBasicInformation `bson:"basicInformation" json:"basicInformation"`
}

type InstanceBasicInformation struct {
	Name       string `bson:"name" json:"name"`
	ExternalID string `bson:"externalId" json:"externalId"`
	Parent     string `bson:"parent" json:"parent"`
}

type InstanceAttribute struct {
}
