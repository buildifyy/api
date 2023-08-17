package template

import (
	"api/pkg/db"
	"api/pkg/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

type MockedDbRepository struct {
	mock.Mock
}

func (m *MockedDbRepository) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedDbRepository) AddOne(data interface{}) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockedDbRepository) GetAllTemplates(filter primitive.D, options *options.FindOptions) ([]models.Template, error) {
	args := m.Called(filter, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Template), args.Error(1)
}

func (m *MockedDbRepository) GetTemplate(filter primitive.D) (*models.Template, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Template), args.Error(1)
}

func (m *MockedDbRepository) GetTypeDropdownValues(collection string) ([]models.Dropdown, error) {
	args := m.Called(collection)
	return args.Get(0).([]models.Dropdown), args.Error(1)
}

func (m *MockedDbRepository) ReplaceTemplate(filter primitive.D, data interface{}) error {
	args := m.Called(filter, data)
	return args.Error(0)
}

func TestNewService(t *testing.T) {
	mockRepository := &MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}
	newService := NewService(mockRepository)

	assert.Equal(t, mockService, newService)
}

func TestService_AddTemplate_Success_CreatesTemplate(t *testing.T) {
	mockRepository := &MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	mockRepository.On("AddOne", mock.AnythingOfType("models.Template")).Return(nil)

	actual := mockService.AddTemplate(models.Template{})
	assert.Equal(t, nil, actual)

	mockRepository.AssertExpectations(t)
}

func TestService_AddTemplate_Fails_ReturnsDuplicateExternalIdError(t *testing.T) {
	mockRepository := &MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expected := db.ErrDuplicateTemplateExternalId

	mockRepository.On("AddOne", mock.AnythingOfType("models.Template")).Return(expected)

	actual := mockService.AddTemplate(models.Template{})
	assert.Equal(t, expected, actual)

	mockRepository.AssertExpectations(t)
}

func TestService_AddTemplate_Fails_ReturnsError(t *testing.T) {
	mockRepository := &MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expected := errors.New("error creating template")

	mockRepository.On("AddOne", mock.AnythingOfType("models.Template")).Return(expected)

	actual := mockService.AddTemplate(models.Template{})
	assert.Equal(t, expected, actual)

	mockRepository.AssertExpectations(t)
}

func TestService_GetTemplate_Success_ReturnsTemplate(t *testing.T) {
	mockRepository := &MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expected := &models.Template{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
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
	mockRepository := &MockedDbRepository{}
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
	mockRepository := &MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expected := make([]models.Template, 0)
	expected = append(expected,
		models.Template{
			TenantID: "the-binary",
			BasicInformation: models.BasicInformation{
				Parent:     "p.com.asset",
				ExternalID: "testtemplate",
				Name:       "Test Template",
				IsCustom:   true,
			},
			Attributes:  nil,
			MetricTypes: nil,
		}, models.Template{
			TenantID: "the-binary",
			BasicInformation: models.BasicInformation{
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
	mockRepository := &MockedDbRepository{}
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
	mockRepository := &MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	mockRepository.On("ReplaceTemplate", mock.AnythingOfType("primitive.D"), mock.AnythingOfType("models.Template")).Return(nil)

	actualErr := mockService.UpdateTemplate("the-binary", models.Template{})
	assert.Nil(t, actualErr)

	mockRepository.AssertExpectations(t)
}

func TestService_UpdateTemplate_Fails_ReturnsError(t *testing.T) {
	mockRepository := &MockedDbRepository{}
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
	mockRepository := &MockedDbRepository{}
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
			BasicInformation: models.BasicInformation{
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
			BasicInformation: models.BasicInformation{
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
	mockRepository := &MockedDbRepository{}
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
