package template

import (
	"api/pkg/db"
	"api/pkg/models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strings"
)

type Controller interface {
	CreateTemplate(c *gin.Context)
	GetParentTemplates(c *gin.Context)
	GetTemplatesList(c *gin.Context)
	GetTemplateById(c *gin.Context)
}

type controller struct {
	templateService Service
}

func NewController(templateService Service) Controller {
	return &controller{
		templateService: templateService,
	}
}

func (c *controller) CreateTemplate(context *gin.Context) {
	tenantID := context.Param("tenantId")
	bytesData, err := io.ReadAll(context.Request.Body)
	if err != nil {
		log.Println("error reading request body: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	var templateToAdd models.Template

	if err := json.Unmarshal(bytesData, &templateToAdd); err != nil {
		log.Println("error parsing request body: ", err)
		context.Status(http.StatusBadRequest)
		return
	}

	templateToAdd.TenantID = tenantID
	templateToAdd.BasicInformation.ExternalID = strings.ToLower(templateToAdd.BasicInformation.ExternalID)

	parentTemplate, err := c.templateService.GetTemplate(tenantID, templateToAdd.BasicInformation.Parent)
	if err != nil {
		log.Println("error fetching parent template: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	templateToAdd.Attributes = append(parentTemplate.Attributes, templateToAdd.Attributes...)
	templateToAdd.MetricTypes = append(parentTemplate.MetricTypes, templateToAdd.MetricTypes...)

	if err := c.templateService.AddTemplate(templateToAdd); err != nil {
		log.Println("error adding template: ", err)
		if errors.Is(err, db.ErrDuplicateTemplateExternalId) {
			context.Status(http.StatusConflict)
			return
		}
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusCreated)
}

func (c *controller) GetParentTemplates(context *gin.Context) {
	tenantID := context.Param("tenantId")

	res, err := c.templateService.GetParentTemplates(tenantID)
	if err != nil {
		log.Println("error getting parent templates: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": res})
}

func (c *controller) GetTemplatesList(context *gin.Context) {
	tenantID := context.Param("tenantId")

	res, err := c.templateService.GetTemplates(tenantID)
	if err != nil {
		log.Println("error getting templates: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": res})
}

func (c *controller) GetTemplateById(context *gin.Context) {
	tenantID := context.Param("tenantId")
	templateID := context.Param("templateId")

	res, err := c.templateService.GetTemplate(tenantID, templateID)
	if err != nil {
		log.Println("error getting templates: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": res})
}
