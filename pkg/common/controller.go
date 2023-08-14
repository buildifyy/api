package common

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Controller interface {
	GetAttributeTypes(c *gin.Context)
	GetMetricTypes(c *gin.Context)
}

type controller struct {
	commonService Service
}

func NewController(commonService Service) Controller {
	return &controller{
		commonService: commonService,
	}
}

func (c *controller) GetAttributeTypes(context *gin.Context) {
	values, err := c.commonService.GetAttributeDropdown()
	if err != nil {
		log.Println("error fetching attribute dropdown values: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": values})
}

func (c *controller) GetMetricTypes(context *gin.Context) {
	values, err := c.commonService.GetMetricTypeDropdown()
	if err != nil {
		log.Println("error fetching metric type dropdown values: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": values})
}
