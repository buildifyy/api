package common

import (
	"api/pkg/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetAttributeDropdown() ([]models.Dropdown, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Dropdown), args.Error(1)
}

func (m *MockService) GetMetricTypeDropdown() ([]models.Dropdown, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Dropdown), args.Error(1)
}

func (m *MockService) GetUnitDropdown() ([]models.Dropdown, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Dropdown), args.Error(1)
}

func (m *MockService) GetRelationships() ([]models.Relationship, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Relationship), args.Error(1)
}

func TestNewController(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		commonService: mockService,
	}

	newController := NewController(mockService)

	assert.Equal(t, mockController, newController)
}

func TestController_GetAttributeTypes_Success_ReturnsOk(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		commonService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Method = "GET"
	ctx.AddParam("tenantId", "the-binary")

	attributeTypesDropdown := []models.Dropdown{
		{
			Label: "Integer",
			Value: "integer",
		},
		{
			Label: "String",
			Value: "string",
		},
	}

	mockService.On("GetAttributeDropdown").Return(attributeTypesDropdown, nil)

	mockController.GetAttributeTypes(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}

func TestController_GetAttributeTypes_Fails_ReturnsInternalServerError(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		commonService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Method = "GET"
	ctx.AddParam("tenantId", "the-binary")

	mockService.On("GetAttributeDropdown").Return(nil, errors.New("error getting attribute dropdown values"))

	mockController.GetAttributeTypes(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}

func TestController_GetMetricTypes_Success_ReturnsOk(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		commonService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Method = "GET"
	ctx.AddParam("tenantId", "the-binary")

	metricTypesDropdown := []models.Dropdown{
		{
			Label: "Integer",
			Value: "integer",
		},
		{
			Label: "String",
			Value: "string",
		},
	}

	mockService.On("GetMetricTypeDropdown").Return(metricTypesDropdown, nil)

	mockController.GetMetricTypes(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}

func TestController_GetMetricTypes_Fails_ReturnsInternalServerError(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		commonService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Method = "GET"
	ctx.AddParam("tenantId", "the-binary")

	mockService.On("GetMetricTypeDropdown").Return(nil, errors.New("error getting metric type dropdown values"))

	mockController.GetMetricTypes(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}
