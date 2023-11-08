package common

import (
	"api/pkg/db"
	"api/pkg/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewService(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	newService := NewService(mockRepository)

	assert.Equal(t, mockService, newService)
}

func TestService_GetAttributeDropdown_Success_ReturnsAttributeTypesDropdown(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expectedDropdownValues := []models.Dropdown{
		{
			Label: "Integer",
			Value: "integer",
		},
		{
			Label: "String",
			Value: "string",
		},
	}
	mockRepository.On("GetTypeDropdownValues", mock.AnythingOfType("string")).Return(expectedDropdownValues, nil)

	actual, actualErr := mockService.GetAttributeDropdown()

	assert.Equal(t, expectedDropdownValues, actual)
	assert.Nil(t, actualErr)

	mockRepository.AssertExpectations(t)
}

func TestService_GetAttributeDropdown_Fails_ReturnsError(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expectedErr := errors.New("error fetching dropdown values")
	mockRepository.On("GetTypeDropdownValues", mock.AnythingOfType("string")).Return(nil, expectedErr)

	actual, actualErr := mockService.GetAttributeDropdown()

	assert.Equal(t, expectedErr, actualErr)
	assert.Nil(t, actual)

	mockRepository.AssertExpectations(t)
}

func TestService_GetMetricTypeDropdown_Success_ReturnsMetricTypesDropdown(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expectedDropdownValues := []models.Dropdown{
		{
			Label: "Integer",
			Value: "integer",
		},
		{
			Label: "String",
			Value: "string",
		},
	}
	mockRepository.On("GetTypeDropdownValues", mock.AnythingOfType("string")).Return(expectedDropdownValues, nil)

	actual, actualErr := mockService.GetMetricTypeDropdown()

	assert.Equal(t, expectedDropdownValues, actual)
	assert.Nil(t, actualErr)

	mockRepository.AssertExpectations(t)
}

func TestService_GetMetricTypeDropdown_Fails_ReturnsError(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expectedErr := errors.New("error fetching dropdown values")
	mockRepository.On("GetTypeDropdownValues", mock.AnythingOfType("string")).Return(nil, expectedErr)

	actual, actualErr := mockService.GetMetricTypeDropdown()

	assert.Equal(t, expectedErr, actualErr)
	assert.Nil(t, actual)

	mockRepository.AssertExpectations(t)
}
