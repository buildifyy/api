package template

import (
	"api/pkg/db"
	"api/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Service interface {
	AddTemplate(template models.Template) error
	GetTemplates(tenantId string) ([]models.Template, error)
	GetTemplate(tenantId string, templateId string) (*models.Template, error)
	GetParentTemplates(tenantId string) ([]models.Dropdown, error)
}

type service struct {
	db db.Repository
}

func NewService(dbRepository db.Repository) Service {
	return &service{
		db: dbRepository,
	}
}

func (s *service) GetParentTemplates(tenantId string) ([]models.Dropdown, error) {
	filter := bson.D{{"tenantId", tenantId}}
	opts := options.Find().SetSort(bson.D{{"basicInformation.name", 1}})
	templates, err := s.db.GetAllTemplates(filter, opts)
	if err != nil {
		log.Println("error fetching all templates: ", err)
		return nil, err
	}

	result := make([]models.Dropdown, 0)
	for _, template := range templates {
		option := models.Dropdown{
			Label: template.BasicInformation.Name,
			Value: template.BasicInformation.ExternalID,
		}
		result = append(result, option)
	}

	return result, nil
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
	opts := options.Find().SetSort(bson.D{{"basicInformation.externalId", 1}})

	templates, err := s.db.GetAllTemplates(filter, opts)
	if err != nil {
		log.Println("error getting all templates: ", err)
		return nil, err
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
