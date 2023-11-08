package common

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, commonController Controller) {
	r.GET("/api/v1/attribute-types", commonController.GetAttributeTypes)
	r.GET("/api/v1/metric-types", commonController.GetMetricTypes)
	r.GET("/api/v1/units", commonController.GetUnits)
	r.GET("/api/v1/relationships", commonController.GetRelationships)
}
