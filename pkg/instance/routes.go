package instance

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, instanceController Controller) {
	r.GET("/api/v1/tenants/:tenantId/instances/form/:parentExternalId", instanceController.GetCreateInstanceForm)
}
