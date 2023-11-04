package template

import (
	"api/pkg/db"
	"api/pkg/models"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	CreateTemplate(c *gin.Context)
	GetParentTemplates(c *gin.Context)
	GetTemplatesList(c *gin.Context)
	GetTemplateById(c *gin.Context)
	UpdateTemplateById(c *gin.Context)
}

type controller struct {
	templateService Service
}

func NewController(templateService Service) Controller {
	return &controller{
		templateService: templateService,
	}
}

func (c *controller) UpdateTemplateById(context *gin.Context) {
	tenantID := context.Param("tenantId")

	var templateToUpdate models.Template
	if err := context.ShouldBindJSON(&templateToUpdate); err != nil {
		log.Println("error parsing request body: ", err)
		context.Status(http.StatusBadRequest)
		return
	}

	if err := c.templateService.UpdateTemplate(tenantID, templateToUpdate); err != nil {
		log.Println("error updating template: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusOK)
}

func (c *controller) CreateTemplate(context *gin.Context) {
	tenantID := context.Param("tenantId")
	var templateToAdd models.Template

	if err := context.ShouldBindJSON(&templateToAdd); err != nil {
		log.Println("error parsing request body: ", err)
		context.Status(http.StatusBadRequest)
		return
	}

	if err := c.templateService.AddTemplate(tenantID, templateToAdd); err != nil {
		log.Println("error adding template: ", err)
		if errors.Is(err, db.ErrDuplicateExternalId) {
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
