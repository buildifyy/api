package common

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetAttributeTypes(c *gin.Context)
	GetMetricTypes(c *gin.Context)
	GetUnits(c *gin.Context)
	GetRelationships(c *gin.Context)
}

type controller struct {
	commonService Service
}

func NewController(commonService Service) Controller {
	return &controller{
		commonService: commonService,
	}
}

func (c *controller) GetRelationships(context *gin.Context) {
	values, err := c.commonService.GetRelationships()
	if err != nil {
		log.Println("error fetching relationships: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": values})
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

func (c *controller) GetUnits(context *gin.Context) {
	values, err := c.commonService.GetUnitDropdown()
	if err != nil {
		log.Println("error fetching unit dropdown values: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": values})
}
