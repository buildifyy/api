package template

import (
	"api/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) AddTemplate(tenantId string, template models.Template) error {
	args := m.Called(tenantId, template)
	return args.Error(0)
}

func (m *MockService) GetTemplates(tenantId string) ([]models.Template, error) {
	args := m.Called(tenantId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Template), args.Error(1)
}

func (m *MockService) GetTemplate(tenantId string, templateId string) (*models.Template, error) {
	args := m.Called(tenantId, templateId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Template), args.Error(1)
}

func (m *MockService) GetParentTemplates(tenantId string) ([]models.Dropdown, error) {
	args := m.Called(tenantId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Dropdown), args.Error(1)
}

func (m *MockService) UpdateTemplate(tenantId string, template models.Template) error {
	args := m.Called(tenantId, template)
	return args.Error(0)
}
