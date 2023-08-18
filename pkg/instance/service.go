package instance

import (
	"api/pkg/db"
	"api/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strings"
)

type Service interface {
	AddInstance(tenantId string, instance models.Instance) error
	GetCreateInstanceForm(tenantId string, parentTemplateExternalId string) error
}

type service struct {
	db db.Repository
}

func NewService(dbRepository db.Repository) Service {
	return &service{
		db: dbRepository,
	}
}

func (s *service) GetCreateInstanceForm(tenantId string, parentTemplateExternalId string) error {
	parentTemplateFilter := bson.D{{"tenantId", tenantId}, {"basicInformation.externalId", parentTemplateExternalId}}
	parentTemplate, err := s.db.GetTemplate(parentTemplateFilter)
	if err != nil {
		log.Println("error getting template: ", err)
		return err
	}

	var ret models.Instance

	parentAttributes := parentTemplate.Attributes

	ret.Attributes = make([]models.InstanceAttribute, 0)
	for _, attr := range parentAttributes {
		if attr.ID != "a25aefe5-b5aa-44b9-9ddf-1f911d1af502" && attr.ID != "c2134cea-ddd2-43f7-a775-e4d12742ef79" {
			ret.Attributes = append(ret.Attributes, models.InstanceAttribute{
				ID:         attr.ID,
				Name:       attr.Name,
				DataType:   attr.DataType,
				IsRequired: attr.IsRequired,
				IsHidden:   attr.IsHidden,
			})
		}
	}

	return nil
}

func (s *service) AddInstance(tenantId string, instance models.Instance) error {
	instance.TenantID = tenantId
	instance.BasicInformation.ExternalID = strings.ToLower(instance.BasicInformation.ExternalID)
	if err := s.db.AddOne("instances", instance); err != nil {
		log.Println("error adding instance: ", err)
		return err
	}
	return nil
}
