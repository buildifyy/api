package template

import (
	"api/pkg/db"
	"api/pkg/models"
	"log"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	AddTemplate(tenantId string, template models.Template) error
	GetTemplates(tenantId string) ([]models.Template, error)
	GetTemplate(tenantId string, templateId string) (*models.Template, error)
	GetParentTemplates(tenantId string) ([]models.ParentTemplateDropdown, error)
	UpdateTemplate(tenantId string, template models.Template) error
}

type service struct {
	db db.Repository
}

func NewService(dbRepository db.Repository) Service {
	return &service{
		db: dbRepository,
	}
}

func (s *service) UpdateTemplate(tenantId string, template models.Template) error {
	template.TenantID = tenantId
	template.BasicInformation.ExternalID = strings.ToLower(template.BasicInformation.ExternalID)
	filter := bson.D{{Key: "tenantId", Value: tenantId}, {Key: "basicInformation.externalId", Value: template.BasicInformation.ExternalID}}

	for i, attribute := range template.Attributes {
		if attribute.ID == "" {
			attributeID, _ := uuid.NewUUID()
			attribute.ID = attributeID.String()
			template.Attributes[i] = attribute
		}
	}

	for i, metric := range template.Metrics {
		if metric.ID == "" {
			metricId, _ := uuid.NewUUID()
			template.Metrics[i].ID = metricId.String()
		}
	}

	if err := s.db.ReplaceTemplate(filter, template); err != nil {
		log.Println("error updating template: ", err)
		return err
	}

	return nil
}

func (s *service) GetParentTemplates(tenantId string) ([]models.ParentTemplateDropdown, error) {
	filter := bson.D{{Key: "tenantId", Value: tenantId}}
	opts := options.Find().SetSort(bson.D{{Key: "basicInformation.name", Value: 1}})
	templates, err := s.db.GetAllTemplates(filter, opts)
	if err != nil {
		log.Println("error fetching all templates: ", err)
		return nil, err
	}

	result := make([]models.ParentTemplateDropdown, 0)
	for _, template := range templates {
		option := models.ParentTemplateDropdown{
			Label:        template.BasicInformation.Name,
			Value:        template.BasicInformation.ExternalID,
			RootTemplate: template.BasicInformation.RootTemplate,
		}
		result = append(result, option)
	}

	return result, nil
}

func (s *service) AddTemplate(tenantId string, template models.Template) error {
	if len(template.Attributes) > 0 {
		for i, attribute := range template.Attributes {
			attributeID, _ := uuid.NewUUID()
			attribute.ID = attributeID.String()
			template.Attributes[i] = attribute
		}
	}
	for i := range template.Metrics {
		metricId, _ := uuid.NewUUID()
		template.Metrics[i].ID = metricId.String()
	}
	template.TenantID = tenantId
	template.BasicInformation.ExternalID = strings.ToLower(template.BasicInformation.ExternalID)

	parentTemplate, err := s.GetTemplate(tenantId, template.BasicInformation.Parent)
	if err != nil {
		log.Println("error fetching parent template: ", err)
		return err
	}
	if parentTemplate.BasicInformation.RootTemplate == "" {
		template.BasicInformation.RootTemplate = parentTemplate.BasicInformation.ExternalID
	} else {
		template.BasicInformation.RootTemplate = parentTemplate.BasicInformation.RootTemplate
	}

	template.Attributes = append(parentTemplate.Attributes, template.Attributes...)
	template.Metrics = append(parentTemplate.Metrics, template.Metrics...)

	if err := s.db.AddOne("templates", template); err != nil {
		log.Println("error inserting template: ", err)
		return err
	}

	return nil
}

func (s *service) GetTemplates(tenantId string) ([]models.Template, error) {
	filter := bson.D{{Key: "tenantId", Value: tenantId}}
	opts := options.Find().SetSort(bson.D{{Key: "basicInformation.externalId", Value: 1}})

	templates, err := s.db.GetAllTemplates(filter, opts)
	if err != nil {
		log.Println("error getting all templates: ", err)
		return nil, err
	}

	return templates, nil
}

func (s *service) GetTemplate(tenantId string, templateId string) (*models.Template, error) {
	filter := bson.D{{Key: "tenantId", Value: tenantId}, {Key: "basicInformation.externalId", Value: templateId}}

	template, err := s.db.GetTemplate(filter)
	if err != nil {
		log.Println("error getting template: ", err)
		return nil, err
	}

	return template, nil
}
