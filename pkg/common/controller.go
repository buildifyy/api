package common

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RegisterRoutes(r *gin.Engine, commonService Service) {
	r.GET("/api/v1/attribute-types", func(c *gin.Context) {
		values, err := commonService.GetAttributeDropdown()
		if err != nil {
			log.Println("error fetching attribute dropdown values: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": values})
	})

	r.GET("/api/v1/metric-types", func(c *gin.Context) {
		values, err := commonService.GetMetricTypeDropdown()
		if err != nil {
			log.Println("error fetching metric type dropdown values: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": values})
	})
}
