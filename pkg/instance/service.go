package instance

import (
	"api/pkg/db"
	"api/pkg/models"
	"api/pkg/template"
	"cmp"
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	AddInstance(tenantId string, instance models.Instance) error
	GetCreateInstanceForm(tenantId string, parentTemplateExternalId string) (*models.InstanceFormMetaData, error)
	GetInstances(tenantId string) ([]models.Instance, error)
	GetInstance(tenantId string, templateId string) (*models.Instance, error)
}

type service struct {
	db              db.Repository
	templateService template.Service
}

func NewService(dbRepository db.Repository, templateService template.Service) Service {
	return &service{
		db:              dbRepository,
		templateService: templateService,
	}
}

func (s *service) GetInstances(tenantId string) ([]models.Instance, error) {
	filter := bson.D{{Key: "tenantId", Value: tenantId}}
	opts := options.Find().SetSort(bson.D{{Key: "basicInformation.externalId", Value: 1}})

	instances, err := s.db.GetAllInstances(filter, opts)
	if err != nil {
		log.Println("error getting all instances: ", err)
		return nil, err
	}

	return instances, nil
}

func (s *service) GetInstance(tenantId string, templateId string) (*models.Instance, error) {
	filter := bson.D{{Key: "tenantId", Value: tenantId}, {Key: "basicInformation.externalId", Value: templateId}}

	template, err := s.db.GetInstance(filter)
	if err != nil {
		log.Println("error getting template: ", err)
		return nil, err
	}

	return template, nil
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
			ID:         "c2134cea-ddd2-43f7-a775-e4d12742ef79",
			Label:      "Name",
			InfoText:   "This will be the name of your instance.",
			Type:       "string",
			TypeLabel:  "String",
			IsRequired: true,
			IsHidden:   false,
		})
	}

	if externalIdAttributeExists := slices.ContainsFunc(parentAttributes, func(attribute models.TemplateAttribute) bool {
		return attribute.ID == "a25aefe5-b5aa-44b9-9ddf-1f911d1af502"
	}); externalIdAttributeExists {
		ret.BasicInformation.Fields = append(ret.BasicInformation.Fields, models.InstanceMetaDataFields{
			ID:         "a25aefe5-b5aa-44b9-9ddf-1f911d1af502",
			Label:      "External ID",
			InfoText:   "A unique identifier for your instance.",
			Type:       "string",
			TypeLabel:  "String",
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

	ret.Metrics.Fields = make([]models.InstanceMetaDataFields, 0)
	for _, metric := range parentTemplate.Metrics {
		dropdownValues := make([]string, 0)
		manualValue := metric.Value
		if metric.IsCalculated {
			dropdownValues = append(dropdownValues, "Calculated")
		}
		if metric.IsSourced {
			dropdownValues = append(dropdownValues, "Sourced")
		}
		if metric.IsManual {
			dropdownValues = append(dropdownValues, "Manual")
		}
		ret.Metrics.Fields = append(ret.Metrics.Fields, models.InstanceMetaDataFields{
			ID:             metric.ID,
			Label:          metric.Name,
			Type:           metric.MetricType,
			DropdownValues: dropdownValues,
			ManualValue:    manualValue,
		})
	}

	return &ret, nil
}

func (s *service) AddInstance(tenantId string, instance models.Instance) error {
	if instance.BasicInformation.Name == "" {
		return fmt.Errorf("%s is required but not provided", "Name")
	}

	if instance.BasicInformation.ExternalId == "" {
		return fmt.Errorf("%s is required but not provided", "External Id")
	}

	if instance.BasicInformation.Parent == "" {
		return fmt.Errorf("%s is required but not provided", "Parent")
	}

	instance.BasicInformation.IsCustom = true

	parentTemplate, err := s.templateService.GetTemplate(tenantId, instance.BasicInformation.Parent)
	if err != nil {
		log.Println("error getting template: ", err)
		return err
	}

	if err := validateAttributes(instance.Attributes, parentTemplate.Attributes); err != nil {
		log.Println("error validating attribute: ", err)
		return err
	}

	if err := validateMetrics(instance.Metrics, parentTemplate.Metrics); err != nil {
		log.Println("error validating metric: ", err)
		return err
	}

	instance.TenantID = tenantId
	instance.BasicInformation.ExternalId = strings.ToLower(instance.BasicInformation.ExternalId)

	if parentTemplate.BasicInformation.RootTemplate == "" {
		instance.BasicInformation.RootTemplate = parentTemplate.BasicInformation.ExternalID
	} else {
		instance.BasicInformation.RootTemplate = parentTemplate.BasicInformation.RootTemplate
	}

	if err := s.db.AddOne("instances", instance); err != nil {
		log.Println("error adding instance: ", err)
		return err
	}
	return nil
}

func validateAttributes(instanceAttributes []models.InstanceAttribute, templateAttributes []models.TemplateAttribute) error {
	for _, attribute := range templateAttributes {
		if attribute.ID == "a25aefe5-b5aa-44b9-9ddf-1f911d1af502" || attribute.ID == "c2134cea-ddd2-43f7-a775-e4d12742ef79" {
			continue
		}
		if attribute.IsRequired {
			if exists := slices.ContainsFunc(instanceAttributes, func(ia models.InstanceAttribute) bool {
				return ia.ID == attribute.ID
			}); !exists {
				return fmt.Errorf("attribute %s is required but not provided", attribute.Name)
			}
		}
	}

	for i, attribute := range instanceAttributes {
		if isValidAttribute := slices.ContainsFunc(templateAttributes, func(ta models.TemplateAttribute) bool {
			attributeId := attribute.ID
			attributeValue := attribute.Value.(string)
			if ta.ID == attribute.ID {
				if ta.IsRequired && len(attributeValue) == 0 {
					log.Printf("attribute %s marked as required is empty", attributeId)
					return false
				}
				if attributeValue != "" {
					switch ta.DataType {
					case "integer":
						integerValue, err := strconv.Atoi(attributeValue)
						if err != nil {
							log.Printf("attribute %s is not an integer value", attributeId)
							return false
						}
						instanceAttributes[i].Value = integerValue
					case "float":
						floatValue, err := strconv.ParseFloat(attributeValue, 64)
						if err != nil {
							log.Printf("attribute %s is not a float value", attributeId)
							return false
						}
						instanceAttributes[i].Value = floatValue
					case "bool":
						booleanValue, err := strconv.ParseBool(strings.ToLower(attributeValue))
						if err != nil {
							log.Printf("attribute %s is not a boolean value", attributeId)
							return false
						}
						instanceAttributes[i].Value = booleanValue
					case "string":
						match, _ := regexp.MatchString("^[a-zA-Z0-9\\s]*$", attributeValue)
						if !match {
							log.Printf("attribute %s is not a valid string", attributeId)
						}
						instanceAttributes[i].Value = attributeValue
					}
				}

				return true
			}

			return false
		}); !isValidAttribute {
			return errors.New("error validating attributes")
		}
		continue
	}
	return nil
}

func validateMetrics(instanceMetrics []models.InstanceMetric, templateMetrics []models.TemplateMetric) error {
	for i, metric := range instanceMetrics {
		if metric.MetricBehaviour == "Manual" {
			if isValidMetricValue := slices.ContainsFunc(templateMetrics, func(tm models.TemplateMetric) bool {
				metricId := metric.ID
				metricValue := metric.Value.(string)
				if tm.ID == metric.ID {
					if metricValue != "" {
						switch metric.MetricType {
						case "integer":
							integerValue, err := strconv.Atoi(metricValue)
							if err != nil {
								log.Printf("metric %s is not an integer value", metricId)
								return false
							}
							instanceMetrics[i].Value = integerValue
						case "float":
							floatValue, err := strconv.ParseFloat(metricValue, 64)
							if err != nil {
								log.Printf("metric %s is not a float value", metricId)
								return false
							}
							instanceMetrics[i].Value = floatValue
						case "bool":
							booleanValue, err := strconv.ParseBool(strings.ToLower(metricValue))
							if err != nil {
								log.Printf("metric %s is not a boolean value", metricId)
								return false
							}
							instanceMetrics[i].Value = booleanValue
						case "string":
							match, _ := regexp.MatchString("^[a-zA-Z0-9\\s]*$", metricValue)
							if !match {
								log.Printf("metric %s is not a valid string", metricId)
							}
							instanceMetrics[i].Value = metricValue
						}
					}

					return true
				}

				return false
			}); !isValidMetricValue {
				return errors.New("error validating metrics")
			}
		}
		continue
	}
	return nil
}
