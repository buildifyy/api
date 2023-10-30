package instance

import (
	"api/pkg/db"
	"api/pkg/models"
	"cmp"
	"log"
	"slices"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	AddInstance(tenantId string, instance models.Instance) error
	GetCreateInstanceForm(tenantId string, parentTemplateExternalId string) (*models.InstanceFormMetaData, error)
}

type service struct {
	db db.Repository
}

func NewService(dbRepository db.Repository) Service {
	return &service{
		db: dbRepository,
	}
}

func (s *service) GetCreateInstanceForm(tenantId string, parentTemplateExternalId string) (*models.InstanceFormMetaData, error) {
	parentTemplateFilter := bson.D{{Key: "tenantId", Value: tenantId}, {Key: "basicInformation.externalId", Value: parentTemplateExternalId}}
	parentTemplate, err := s.db.GetTemplate(parentTemplateFilter)
	if err != nil {
		log.Println("error getting template: ", err)
		return nil, err
	}

	attributeTypes, err := s.db.GetTypeDropdownValues("attribute_types")
	if err != nil {
		log.Println("error finding attribute types: ", err)
		return nil, err
	}
	slices.SortFunc(attributeTypes, func(a, b models.Dropdown) int {
		return cmp.Compare(strings.ToLower(a.Label), strings.ToLower(b.Label))
	})

	var ret models.InstanceFormMetaData

	parentAttributes := parentTemplate.Attributes

	ret.BasicInformation.Fields = make([]models.InstanceMetaDataFields, 0)
	if nameAttributeExists := slices.ContainsFunc(parentAttributes, func(attribute models.TemplateAttribute) bool {
		return attribute.ID == "c2134cea-ddd2-43f7-a775-e4d12742ef79"
	}); nameAttributeExists {
		ret.BasicInformation.Fields = append(ret.BasicInformation.Fields, models.InstanceMetaDataFields{
			Label:      "Name",
			InfoText:   "This will be the name of your instance.",
			Type:       "string",
			IsRequired: true,
			IsHidden:   false,
		})
	}

	if externalIdAttributeExists := slices.ContainsFunc(parentAttributes, func(attribute models.TemplateAttribute) bool {
		return attribute.ID == "a25aefe5-b5aa-44b9-9ddf-1f911d1af502"
	}); externalIdAttributeExists {
		ret.BasicInformation.Fields = append(ret.BasicInformation.Fields, models.InstanceMetaDataFields{
			Label:      "External ID",
			InfoText:   "A unique identifier for your instance.",
			Type:       "string",
			IsRequired: true,
			IsHidden:   false,
		})
	}

	ret.Attributes.Fields = make([]models.InstanceMetaDataFields, 0)
	for _, attr := range parentAttributes {
		if attr.ID != "a25aefe5-b5aa-44b9-9ddf-1f911d1af502" && attr.ID != "c2134cea-ddd2-43f7-a775-e4d12742ef79" {
			attributeTypeIndex, _ := slices.BinarySearchFunc(attributeTypes, models.Dropdown{
				Value: attr.DataType,
			}, func(dropdown models.Dropdown, dropdown2 models.Dropdown) int {
				return cmp.Compare(dropdown.Value, dropdown2.Value)
			})
			ret.Attributes.Fields = append(ret.Attributes.Fields, models.InstanceMetaDataFields{
				ID:         attr.ID,
				Label:      attr.Name,
				TypeLabel:  attributeTypes[attributeTypeIndex].Label,
				Type:       attr.DataType,
				InfoText:   "",
				IsRequired: attr.IsRequired,
				IsHidden:   attr.IsHidden,
			})
		}
	}

	return &ret, nil
}

func (s *service) AddInstance(tenantId string, instance models.Instance) error {
	//instance.TenantID = tenantId
	//instance.BasicInformation.ExternalID = strings.ToLower(instance.BasicInformation.ExternalID)
	//if err := s.db.AddOne("instances", instance); err != nil {
	//	log.Println("error adding instance: ", err)
	//	return err
	//}
	return nil
}
