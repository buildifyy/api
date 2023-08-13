package template

import (
	"api/pkg/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func RegisterRoutes(r *gin.Engine, templateService Service) {
	r.POST("/api/v1/tenants/:tenantId/templates", func(c *gin.Context) {
		tenantID := c.Param("tenantId")
		bytesData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("error reading request body: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		var templateToAdd models.Template

		if err := json.Unmarshal(bytesData, &templateToAdd); err != nil {
			log.Println("error parsing request body: ", err)
			c.Status(http.StatusBadRequest)
			return
		}

		templateToAdd.TenantID = tenantID

		if err := templateService.AddTemplate(templateToAdd); err != nil {
			log.Println("error adding template: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusCreated)
	})

	r.GET("/api/v1/tenants/:tenantId/templates/parent", func(c *gin.Context) {
		tenantID := c.Param("tenantId")

		res, err := templateService.GetParentTemplates(tenantID)
		if err != nil {
			log.Println("error getting parent templates: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
	})

	r.GET("/api/v1/tenants/:tenantId/templates", func(c *gin.Context) {
		tenantID := c.Param("tenantId")

		res, err := templateService.GetTemplates(tenantID)
		if err != nil {
			log.Println("error getting templates: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
	})

	r.GET("/api/v1/tenants/:tenantId/templates/:templateId", func(c *gin.Context) {
		tenantID := c.Param("tenantId")
		templateID := c.Param("templateId")

		res, err := templateService.GetTemplate(tenantID, templateID)
		if err != nil {
			log.Println("error getting templates: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": res})
	})
}
