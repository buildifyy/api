package instance

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, instanceController Controller) {
	r.GET("/api/v1/tenants/:tenantId/instances/form/:parentExternalId", instanceController.GetCreateInstanceForm)
	r.POST("/api/v1/tenants/:tenantId/instances", instanceController.AddInstance)
	r.GET("/api/v1/tenants/:tenantId/instances", instanceController.GetInstanceList)
	r.GET("/api/v1/tenants/:tenantId/instances/:instanceId", instanceController.GetInstanceById)
}
