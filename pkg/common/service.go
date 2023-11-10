package common

import (
	"api/pkg/db"
	"api/pkg/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	GetAttributeDropdown() ([]models.Dropdown, error)
	GetMetricTypeDropdown() ([]models.Dropdown, error)
	GetUnitDropdown() ([]models.Dropdown, error)
	GetRelationships() ([]models.Relationship, error)
	GetRelationship(id string) (models.Relationship, error)
}

type service struct {
	db db.Repository
}

func NewService(repository db.Repository) Service {
	return &service{
		db: repository,
	}
}

func (s *service) GetRelationships() ([]models.Relationship, error) {
	values, err := s.db.GetRelationships(nil, "relationships")
	if err != nil {
		log.Println("error fetching relationships: ", err)
		return nil, err
	}

	return values, nil
}

var objectIDFromHex = func(hex string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		log.Println("error creating object id from hex")
	}
	return objectID
}

func (s *service) GetRelationship(id string) (models.Relationship, error) {
	filter := bson.D{{Key: "_id", Value: objectIDFromHex(id)}}
	values, err := s.db.GetRelationships(filter, "relationships")
	if err != nil {
		log.Println("error fetching relationships: ", err)
		return models.Relationship{}, err
	}

	return values[0], nil
}

func (s *service) GetAttributeDropdown() ([]models.Dropdown, error) {
	values, err := s.db.GetTypeDropdownValues("attribute_types")
	if err != nil {
		log.Println("error fetching dropdown values for attributes: ", err)
		return nil, err
	}

	return values, nil
}

func (s *service) GetMetricTypeDropdown() ([]models.Dropdown, error) {
	values, err := s.db.GetTypeDropdownValues("metric_types")
	if err != nil {
		log.Println("error fetching dropdown values for metric types: ", err)
		return nil, err
	}

	return values, nil
}

func (s *service) GetUnitDropdown() ([]models.Dropdown, error) {
	values, err := s.db.GetTypeDropdownValues("units")
	if err != nil {
		log.Println("error fetching dropdown values for units: ", err)
		return nil, err
	}

	return values, nil
}
