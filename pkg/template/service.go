package template

import (
	"api/pkg/db"
	"api/pkg/models"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

type Service interface {
	AddTemplate(template models.Template) error
	GetTemplates(tenantId string) ([]models.Template, error)
	GetTemplate(tenantId string, templateId string) (*models.Template, error)
}

type service struct {
	db db.Repository
}

func NewService(dbRepository db.Repository) Service {
	return &service{
		db: dbRepository,
	}
}

func (s *service) AddTemplate(template models.Template) error {
	if err := s.db.AddOne(template); err != nil {
		log.Println("error inserting template: ", err)
		return err
	}

	return nil
}

func (s *service) GetTemplates(tenantId string) ([]models.Template, error) {
	filter := bson.D{{"tenantId", tenantId}}

	byteArray, err := s.db.GetAllTemplates(filter)
	if err != nil {
		log.Println("error getting all templates: ", err)
		return nil, err
	}

	templates := make([]models.Template, 0)
	for _, data := range byteArray {
		var template models.Template
		if err := json.Unmarshal(data, &template); err != nil {
			log.Println("error unmarshalling bytes to templates: ", err)
			return nil, err
		}

		templates = append(templates, template)
	}

	return templates, nil
}

func (s *service) GetTemplate(tenantId string, templateId string) (*models.Template, error) {
	filter := bson.D{{"tenantId", tenantId}, {"basicInformation.externalId", templateId}}

	template, err := s.db.GetTemplate(filter)
	if err != nil {
		log.Println("error getting all templates: ", err)
		return nil, err
	}

	return template, nil
}
