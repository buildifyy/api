package template

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, templateController Controller) {
	r.GET("/api/v1/tenants/:tenantId/templates", templateController.GetTemplatesList)
	r.GET("/api/v1/tenants/:tenantId/templates/:templateId", templateController.GetTemplateById)
	r.GET("/api/v1/tenants/:tenantId/templates/parent", templateController.GetParentTemplates)
	r.POST("/api/v1/tenants/:tenantId/templates", templateController.CreateTemplate)
	r.PUT("/api/v1/tenants/:tenantId/templates/:templateId", templateController.UpdateTemplateById)
}
