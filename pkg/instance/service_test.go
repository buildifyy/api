package instance

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

func TestService_AddInstance_Success_CreatesInstance(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	mockRepository.On("AddOne", mock.AnythingOfType("string"), mock.AnythingOfType("models.Instance")).Return(nil)

	actualErr := mockService.AddInstance("the-binary", models.Instance{})

	assert.Nil(t, actualErr)

	mockRepository.AssertExpectations(t)
}

func TestService_AddInstance_Fails_ReturnsError(t *testing.T) {
	mockRepository := &db.MockedDbRepository{}
	mockService := &service{
		db: mockRepository,
	}

	expectedErr := errors.New("error adding instance")
	mockRepository.On("AddOne", mock.AnythingOfType("string"), mock.AnythingOfType("models.Instance")).Return(expectedErr)

	actualErr := mockService.AddInstance("the-binary", models.Instance{})

	assert.Equal(t, expectedErr, actualErr)

	mockRepository.AssertExpectations(t)
}
