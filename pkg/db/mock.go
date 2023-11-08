package db

import (
	"api/pkg/models"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockedDbRepository struct {
	mock.Mock
}

func (m *MockedDbRepository) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedDbRepository) AddOne(collectionName string, data interface{}) error {
	args := m.Called(collectionName, data)
	return args.Error(0)
}

func (m *MockedDbRepository) GetAllTemplates(filter primitive.D, options *options.FindOptions) ([]models.Template, error) {
	args := m.Called(filter, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Template), args.Error(1)
}

func (m *MockedDbRepository) GetAllInstances(filter primitive.D, options *options.FindOptions) ([]models.Instance, error) {
	args := m.Called(filter, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Instance), args.Error(1)
}

func (m *MockedDbRepository) GetTemplate(filter primitive.D) (*models.Template, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Template), args.Error(1)
}

func (m *MockedDbRepository) GetInstance(filter primitive.D) (*models.Instance, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Instance), args.Error(1)
}

func (m *MockedDbRepository) GetTypeDropdownValues(collection string) ([]models.Dropdown, error) {
	args := m.Called(collection)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Dropdown), args.Error(1)
}

func (m *MockedDbRepository) GetRelationships(collection string) ([]models.Relationship, error) {
	args := m.Called(collection)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Relationship), args.Error(1)
}

func (m *MockedDbRepository) ReplaceTemplate(filter primitive.D, data interface{}) error {
	args := m.Called(filter, data)
	return args.Error(0)
}
