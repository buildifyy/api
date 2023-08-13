package common

import (
	"api/pkg/db"
	"api/pkg/models"
	"log"
)

type Service interface {
	GetAttributeDropdown() ([]models.Dropdown, error)
	GetMetricTypeDropdown() ([]models.Dropdown, error)
}

type service struct {
	db db.Repository
}

func NewService(repository db.Repository) Service {
	return &service{
		db: repository,
	}
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
