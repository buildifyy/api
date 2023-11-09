package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Relationship struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Source      string             `bson:"source" json:"source"`
	Target      []string           `bson:"target" json:"target"`
	Cardinality string             `bson:"cardinality" json:"cardinality"`
	Inverse     primitive.ObjectID `bson:"inverse" json:"inverse"`
}
