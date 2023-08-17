package template

import (
	"api/pkg/db"
	"api/pkg/models"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) AddTemplate(template models.Template) error {
	args := m.Called(template)
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

func TestNewController(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	newController := NewController(mockService)

	assert.Equal(t, mockController, newController)
}

func TestController_CreateTemplate_Success(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Method = "POST"
	ctx.AddParam("tenantId", "the-binary")
	jsonBytes, err := json.Marshal(models.Template{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "p.com.asset",
			Name:       "Test Template 1",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	})
	if err != nil {
		panic(err)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	parentTemplateResponse := &models.Template{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "",
			Name:       "Asset",
			ExternalID: "p.com.asset",
			IsCustom:   false,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	}

	mockService.On("GetTemplate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(parentTemplateResponse, nil)
	mockService.On("AddTemplate", mock.AnythingOfType("models.Template")).Return(nil)

	mockController.CreateTemplate(ctx)

	assert.Equal(t, http.StatusCreated, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}

func TestController_CreateTemplate_FailsToParseRequestBody_ReturnsBadRequest(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Method = "POST"
	ctx.AddParam("tenantId", "the-binary")

	//bad json request to create a template
	jsonBytes, err := json.Marshal([]models.Template{{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "p.com.asset",
			Name:       "Test Template 1",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	}})
	if err != nil {
		panic(err)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

	mockController.CreateTemplate(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
	mockService.AssertExpectations(t)
}

func TestController_CreateTemplate_FailsToGetParentTemplateDetails_ReturnsInternalServerError(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Method = "POST"
	ctx.AddParam("tenantId", "the-binary")

	//bad json request to create a template
	jsonBytes, err := json.Marshal(models.Template{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "p.com.asset",
			Name:       "Test Template 1",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	})
	if err != nil {
		panic(err)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

	expectedErr := errors.New("error getting parent template")
	mockService.On("GetTemplate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, expectedErr)

	mockController.CreateTemplate(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	mockService.AssertExpectations(t)
}

func TestController_CreateTemplate_FailsToCreateTemplate_ReturnsInternalServerError(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Method = "POST"
	ctx.AddParam("tenantId", "the-binary")

	//bad json request to create a template
	jsonBytes, err := json.Marshal(models.Template{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "p.com.asset",
			Name:       "Test Template 1",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	})
	if err != nil {
		panic(err)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	parentTemplateResponse := &models.Template{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "",
			Name:       "Asset",
			ExternalID: "p.com.asset",
			IsCustom:   false,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	}

	expectedErr := errors.New("error adding new template")
	mockService.On("GetTemplate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(parentTemplateResponse, nil)
	mockService.On("AddTemplate", mock.AnythingOfType("models.Template")).Return(expectedErr)

	mockController.CreateTemplate(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	mockService.AssertExpectations(t)
}

func TestController_CreateTemplate_FailsToCreateTemplateWhenTemplateAlreadyExists_ReturnsStatusConflict(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Method = "POST"
	ctx.AddParam("tenantId", "the-binary")

	//bad json request to create a template
	jsonBytes, err := json.Marshal(models.Template{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "p.com.asset",
			Name:       "Test Template 1",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	})
	if err != nil {
		panic(err)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	parentTemplateResponse := &models.Template{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "",
			Name:       "Asset",
			ExternalID: "p.com.asset",
			IsCustom:   false,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	}

	expectedErr := db.ErrDuplicateTemplateExternalId
	mockService.On("GetTemplate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(parentTemplateResponse, nil)
	mockService.On("AddTemplate", mock.AnythingOfType("models.Template")).Return(expectedErr)

	mockController.CreateTemplate(ctx)

	assert.Equal(t, http.StatusConflict, ctx.Writer.Status())
	mockService.AssertExpectations(t)
}

func TestController_UpdateTemplateById_Success(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Method = "PUT"
	ctx.AddParam("tenantId", "the-binary")
	ctx.AddParam("templateId", "testtemplate1")

	jsonBytes, err := json.Marshal(models.Template{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "p.com.asset",
			Name:       "Test Template 1",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	})
	if err != nil {
		panic(err)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

	mockService.On("UpdateTemplate", mock.AnythingOfType("string"), mock.AnythingOfType("models.Template")).Return(nil)

	mockController.UpdateTemplateById(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}

func TestController_UpdateTemplateById_FailsToParseRequestBody_ReturnsBadRequest(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Method = "PUT"
	ctx.AddParam("tenantId", "the-binary")
	ctx.AddParam("templateId", "testtemplate1")

	//bad json request to create a template
	jsonBytes, err := json.Marshal([]models.Template{{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "p.com.asset",
			Name:       "Test Template 1",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	}})
	if err != nil {
		panic(err)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

	mockController.UpdateTemplateById(ctx)

	assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
	mockService.AssertExpectations(t)
}

func TestController_UpdateTemplateById_FailsToUpdateTemplate_ReturnsInternalServerError(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Method = "PUT"
	ctx.AddParam("tenantId", "the-binary")
	ctx.AddParam("templateId", "testtemplate1")

	//bad json request to create a template
	jsonBytes, err := json.Marshal(models.Template{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "p.com.asset",
			Name:       "Test Template 1",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	})
	if err != nil {
		panic(err)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

	expectedErr := errors.New("error updating new template")
	mockService.On("UpdateTemplate", mock.AnythingOfType("string"), mock.AnythingOfType("models.Template")).Return(expectedErr)

	mockController.UpdateTemplateById(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	mockService.AssertExpectations(t)
}

func TestController_GetParentTemplates_Success(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Method = "GET"
	ctx.AddParam("tenantId", "the-binary")

	parentTemplatesDropdown := []models.Dropdown{
		{
			Label: "Test Template 1",
			Value: "testtemplate1",
		},
		{
			Label: "Test Template 2",
			Value: "testtemplate2",
		},
	}

	mockService.On("GetParentTemplates", mock.AnythingOfType("string")).Return(parentTemplatesDropdown, nil)

	mockController.GetParentTemplates(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}

func TestController_GetParentTemplates_Fails_ReturnsInternalServerError(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Method = "GET"
	ctx.AddParam("tenantId", "the-binary")

	mockService.On("GetParentTemplates", mock.AnythingOfType("string")).Return(nil, errors.New("error fetching parent templates dropdown"))

	mockController.GetParentTemplates(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}

func TestController_GetTemplatesList_Success(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Method = "GET"
	ctx.AddParam("tenantId", "the-binary")

	templatesResponse := []models.Template{{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "p.com.asset",
			Name:       "Test Template 1",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	}}

	mockService.On("GetTemplates", mock.AnythingOfType("string")).Return(templatesResponse, nil)

	mockController.GetTemplatesList(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}

func TestController_GetTemplatesList_Fails_ReturnsInternalServerError(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Method = "GET"
	ctx.AddParam("tenantId", "the-binary")

	mockService.On("GetTemplates", mock.AnythingOfType("string")).Return(nil, errors.New("error fetching templates"))

	mockController.GetTemplatesList(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}

func TestController_GetTemplatesById_Success(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Method = "GET"
	ctx.AddParam("tenantId", "the-binary")
	ctx.AddParam("templateId", "testtemplate1")

	templateResponse := &models.Template{
		TenantID: "the-binary",
		BasicInformation: models.BasicInformation{
			Parent:     "p.com.asset",
			Name:       "Test Template 1",
			ExternalID: "testtemplate1",
			IsCustom:   true,
		},
		Attributes:  make([]models.Attribute, 0),
		MetricTypes: make([]models.MetricType, 0),
	}

	mockService.On("GetTemplate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(templateResponse, nil)

	mockController.GetTemplateById(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}

func TestController_GetTemplateById_Fails_ReturnsInternalServerError(t *testing.T) {
	mockService := &MockService{}
	mockController := &controller{
		templateService: mockService,
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}
	ctx.Request.Method = "GET"
	ctx.AddParam("tenantId", "the-binary")
	ctx.AddParam("templateId", "testtemplate1")

	mockService.On("GetTemplate", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("error fetching template"))

	mockController.GetTemplateById(ctx)

	assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())

	mockService.AssertExpectations(t)
}
