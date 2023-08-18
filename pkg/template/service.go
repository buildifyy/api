package template

import (
	"api/pkg/db"
	"api/pkg/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
)

type Service interface {
	AddTemplate(tenantId string, template models.Template) error
	GetTemplates(tenantId string) ([]models.Template, error)
	GetTemplate(tenantId string, templateId string) (*models.Template, error)
	GetParentTemplates(tenantId string) ([]models.Dropdown, error)
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
	filter := bson.D{{"tenantId", tenantId}, {"basicInformation.externalId", template.BasicInformation.ExternalID}}

	for i, attribute := range template.Attributes {
		if attribute.ID == "" {
			attributeID, _ := uuid.NewUUID()
			attribute.ID = attributeID.String()
			template.Attributes[i] = attribute
		}
	}

	for i, metricType := range template.MetricTypes {
		if metricType.ID == "" {
			metricTypeID, _ := uuid.NewUUID()
			metricType.ID = metricTypeID.String()
			for j, metric := range metricType.Metrics {
				if metric.ID == "" {
					metricID, _ := uuid.NewUUID()
					metric.ID = metricID.String()
					metricType.Metrics[j] = metric
				}
			}
			template.MetricTypes[i] = metricType
		}
	}

	if err := s.db.ReplaceTemplate(filter, template); err != nil {
		log.Println("error updating template: ", err)
		return err
	}

	return nil
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

func (s *service) AddTemplate(tenantId string, template models.Template) error {
	if len(template.Attributes) > 0 {
		for i, attribute := range template.Attributes {
			attributeID, _ := uuid.NewUUID()
			attribute.ID = attributeID.String()
			template.Attributes[i] = attribute
		}
	}
	if len(template.MetricTypes) > 0 {
		for i, metricType := range template.MetricTypes {
			metricTypeID, _ := uuid.NewUUID()
			metricType.ID = metricTypeID.String()
			for j, metric := range metricType.Metrics {
				metricID, _ := uuid.NewUUID()
				metric.ID = metricID.String()
				metricType.Metrics[j] = metric
			}
			template.MetricTypes[i] = metricType
		}
	}
	template.TenantID = tenantId
	template.BasicInformation.ExternalID = strings.ToLower(template.BasicInformation.ExternalID)

	parentTemplate, err := s.GetTemplate(tenantId, template.BasicInformation.Parent)
	if err != nil {
		log.Println("error fetching parent template: ", err)
		return err
	}

	template.Attributes = append(parentTemplate.Attributes, template.Attributes...)
	template.MetricTypes = append(parentTemplate.MetricTypes, template.MetricTypes...)

	if err := s.db.AddOne("templates", template); err != nil {
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
		log.Println("error getting template: ", err)
		return nil, err
	}

	return template, nil
}
