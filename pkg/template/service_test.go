package template

import (
	"api/pkg/db"
	"api/pkg/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestNewService(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}
	newService := NewService(mockRepository)

	assert.Equal(t, mockService, newService)
}

func TestService_AddTemplate_Success_CreatesTemplate(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	mockRepository.On("GetTemplate", mock.AnythingOfType("primitive.D")).Return(&models.Template{
		TenantID: "the-binary",
		BasicInformation: models.TemplateBasicInformation{
			Parent:     "",
			Name:       "Asset",
			ExternalID: "p.com.asset",
			IsCustom:   false,
		},
		Attributes:  make([]models.TemplateAttribute, 0),
		MetricTypes: make([]models.TemplateMetricType, 0),
	}, nil)
	mockRepository.On("AddOne", mock.AnythingOfType("models.Template")).Return(nil)

	actual := mockService.AddTemplate("the-binary", models.Template{
		TenantID: "the-binary",
		BasicInformation: models.TemplateBasicInformation{
			Name:       "testtemplate1",
			Parent:     "p.com.asset",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes: []models.TemplateAttribute{
			{
				Name:           "attribute1",
				DataType:       "integer",
				IsRequired:     false,
				IsHidden:       false,
				OwningTemplate: "testtemplate1",
			},
		},
		MetricTypes: []models.TemplateMetricType{
			{
				Name:       "metrictype1",
				MetricType: "integer",
				Metrics: []models.TemplateMetric{
					{
						Name:         "metric1",
						IsManual:     false,
						Value:        nil,
						IsCalculated: false,
						IsSourced:    false,
					},
				},
				OwningTemplate: "testtemplate1",
			},
		},
	})
	assert.Equal(t, nil, actual)

	mockRepository.AssertExpectations(t)
}

func TestService_AddTemplate_Fails_ReturnsDuplicateExternalIdError(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expected := db.ErrDuplicateTemplateExternalId

	mockRepository.On("GetTemplate", mock.AnythingOfType("primitive.D")).Return(&models.Template{
		TenantID: "the-binary",
		BasicInformation: models.TemplateBasicInformation{
			Parent:     "",
			Name:       "Asset",
			ExternalID: "p.com.asset",
			IsCustom:   false,
		},
		Attributes:  make([]models.TemplateAttribute, 0),
		MetricTypes: make([]models.TemplateMetricType, 0),
	}, nil)
	mockRepository.On("AddOne", mock.AnythingOfType("models.Template")).Return(expected)

	actual := mockService.AddTemplate("the-binary", models.Template{})
	assert.Equal(t, expected, actual)

	mockRepository.AssertExpectations(t)
}

func TestService_AddTemplate_FailsGettingParentTemplate_ReturnsError(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expectedErr := errors.New("error getting parent template")

	mockRepository.On("GetTemplate", mock.AnythingOfType("primitive.D")).Return(nil, expectedErr)

	actualErr := mockService.AddTemplate("the-binary", models.Template{})
	assert.Equal(t, expectedErr, actualErr)

	mockRepository.AssertExpectations(t)
}

func TestService_AddTemplate_Fails_ReturnsError(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expected := errors.New("error creating template")

	mockRepository.On("GetTemplate", mock.AnythingOfType("primitive.D")).Return(&models.Template{
		TenantID: "the-binary",
		BasicInformation: models.TemplateBasicInformation{
			Parent:     "",
			Name:       "Asset",
			ExternalID: "p.com.asset",
			IsCustom:   false,
		},
		Attributes:  make([]models.TemplateAttribute, 0),
		MetricTypes: make([]models.TemplateMetricType, 0),
	}, nil)
	mockRepository.On("AddOne", mock.AnythingOfType("models.Template")).Return(expected)

	actual := mockService.AddTemplate("the-binary", models.Template{})
	assert.Equal(t, expected, actual)

	mockRepository.AssertExpectations(t)
}

func TestService_GetTemplate_Success_ReturnsTemplate(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expected := &models.Template{
		TenantID: "the-binary",
		BasicInformation: models.TemplateBasicInformation{
			Parent:     "p.com.asset",
			ExternalID: "testtemplate",
			Name:       "Test Template",
			IsCustom:   true,
		},
		Attributes:  nil,
		MetricTypes: nil,
	}

	mockRepository.On("GetTemplate", mock.AnythingOfType("primitive.D")).Return(expected, nil)

	actual, actualErr := mockService.GetTemplate("the-binary", "testtemplate")
	assert.Equal(t, expected, actual)
	assert.Nil(t, actualErr)

	mockRepository.AssertExpectations(t)
}

func TestService_GetTemplate_Fails_ReturnsError(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expectedErr := errors.New("error getting template")

	mockRepository.On("GetTemplate", mock.AnythingOfType("primitive.D")).Return(nil, expectedErr)

	actual, actualErr := mockService.GetTemplate("the-binary", "testtemplate")
	assert.Equal(t, expectedErr, actualErr)
	assert.Nil(t, actual)

	mockRepository.AssertExpectations(t)
}

func TestService_GetTemplates_Success_ReturnsTemplates(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expected := make([]models.Template, 0)
	expected = append(expected,
		models.Template{
			TenantID: "the-binary",
			BasicInformation: models.TemplateBasicInformation{
				Parent:     "p.com.asset",
				ExternalID: "testtemplate",
				Name:       "Test Template",
				IsCustom:   true,
			},
			Attributes:  nil,
			MetricTypes: nil,
		}, models.Template{
			TenantID: "the-binary",
			BasicInformation: models.TemplateBasicInformation{
				Parent:     "p.com.asset",
				ExternalID: "testtemplate2",
				Name:       "Test Template 2",
				IsCustom:   true,
			},
			Attributes:  nil,
			MetricTypes: nil,
		})

	mockRepository.On("GetAllTemplates", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("*options.FindOptions")).Return(expected, nil)

	actual, actualErr := mockService.GetTemplates("the-binary")
	assert.Equal(t, expected, actual)
	assert.Nil(t, actualErr)

	mockRepository.AssertExpectations(t)
}

func TestService_GetTemplates_Fails_ReturnsError(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expectedErr := errors.New("error fetching templates")

	mockRepository.On("GetAllTemplates", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("*options.FindOptions")).Return(nil, expectedErr)

	actual, actualErr := mockService.GetTemplates("the-binary")
	assert.Equal(t, expectedErr, actualErr)
	assert.Nil(t, actual)

	mockRepository.AssertExpectations(t)
}

func TestService_UpdateTemplate_Success(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	mockRepository.On("ReplaceTemplate", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("models.Template")).Return(nil)

	actualErr := mockService.UpdateTemplate("the-binary", models.Template{
		TenantID: "the-binary",
		BasicInformation: models.TemplateBasicInformation{
			Name:       "testtemplate1",
			Parent:     "p.com.asset",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes: []models.TemplateAttribute{
			{
				ID:             "412ba829-eca5-4513-97e7-f30c34f03a70",
				Name:           "attribute1",
				DataType:       "integer",
				IsRequired:     false,
				IsHidden:       false,
				OwningTemplate: "testtemplate1",
			},
			{
				Name:           "attribute2",
				DataType:       "integer",
				IsRequired:     false,
				IsHidden:       false,
				OwningTemplate: "testtemplate1",
			},
		},
		MetricTypes: []models.TemplateMetricType{
			{
				ID:         "549fe288-63ee-480a-9a1d-85d4021b77ca",
				Name:       "metrictype1",
				MetricType: "integer",
				Metrics: []models.TemplateMetric{
					{
						ID:           "90f336bb-1614-42ab-800b-1b1f821161c1",
						Name:         "metric1",
						IsManual:     false,
						Value:        nil,
						IsCalculated: false,
						IsSourced:    false,
					},
				},
				OwningTemplate: "testtemplate1",
			},
			{
				Name:       "metrictype2",
				MetricType: "integer",
				Metrics: []models.TemplateMetric{
					{
						Name:         "metric1",
						IsManual:     false,
						Value:        nil,
						IsCalculated: false,
						IsSourced:    false,
					},
				},
				OwningTemplate: "testtemplate1",
			},
		},
	})
	assert.Nil(t, actualErr)

	mockRepository.AssertExpectations(t)
}

func TestService_UpdateTemplate_Fails_ReturnsError(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expectedErr := errors.New("error replacing template")

	mockRepository.On("ReplaceTemplate", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("models.Template")).Return(expectedErr)

	actualErr := mockService.UpdateTemplate("the-binary", models.Template{})
	assert.Equal(t, expectedErr, actualErr)

	mockRepository.AssertExpectations(t)
}

func TestService_GetParentTemplates_Success_ReturnsExternalIdSlice(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expected := []models.Dropdown{
		{
			Label: "Test Template 1",
			Value: "testtemplate1",
		},
		{
			Label: "Test Template 2",
			Value: "testtemplate2",
		},
	}

	templates := []models.Template{
		{
			TenantID: "the-binary",
			BasicInformation: models.TemplateBasicInformation{
				Parent:     "p.com.asset",
				ExternalID: "testtemplate1",
				Name:       "Test Template 1",
				IsCustom:   true,
			},
			Attributes:  nil,
			MetricTypes: nil,
		},
		{
			TenantID: "the-binary",
			BasicInformation: models.TemplateBasicInformation{
				Parent:     "p.com.asset",
				ExternalID: "testtemplate2",
				Name:       "Test Template 2",
				IsCustom:   true,
			},
			Attributes:  nil,
			MetricTypes: nil,
		},
	}

	mockRepository.On("GetAllTemplates", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("*options.FindOptions")).Return(templates, nil)

	actual, actualErr := mockService.GetParentTemplates("the-binary")
	assert.Equal(t, expected, actual)
	assert.Nil(t, actualErr)

	mockRepository.AssertExpectations(t)
}

func TestService_GetParentTemplates_Fails_ReturnsError(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expectedErr := errors.New("error fetching parent templates")

	mockRepository.On("GetAllTemplates", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("*options.FindOptions")).Return(nil, expectedErr)

	actual, actualErr := mockService.GetParentTemplates("the-binary")
	assert.Equal(t, expectedErr, actualErr)
	assert.Nil(t, actual)

	mockRepository.AssertExpectations(t)
}
