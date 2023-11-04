package instance

import (
	"api/pkg/db"
	"api/pkg/models"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	AddInstance(c *gin.Context)
	GetCreateInstanceForm(c *gin.Context)
	GetInstanceList(c *gin.Context)
}

type controller struct {
	instanceService Service
}

func NewController(instanceService Service) Controller {
	return &controller{
		instanceService: instanceService,
	}
}

func (c *controller) AddInstance(context *gin.Context) {
	tenantId := context.Param("tenantId")
	var instanceToAdd models.Instance

	if err := context.ShouldBindJSON(&instanceToAdd); err != nil {
		log.Println("error parsing request body: ", err)
		context.Status(http.StatusBadRequest)
		return
	}

	if err := c.instanceService.AddInstance(tenantId, instanceToAdd); err != nil {
		log.Println("error adding instance: ", err)
		if strings.Contains(err.Error(), "error validating attribute") || strings.Contains(err.Error(), "is required but not provided") {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, db.ErrDuplicateExternalId) {
			context.Status(http.StatusConflict)
			return
		}
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusCreated)
}

func (c *controller) GetInstanceList(context *gin.Context) {
	tenantID := context.Param("tenantId")

	res, err := c.instanceService.GetInstances(tenantID)
	if err != nil {
		log.Println("error getting instances: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": res})
}

func (c *controller) GetCreateInstanceForm(context *gin.Context) {
	tenantId := context.Param("tenantId")
	parentExternalId := context.Param("parentExternalId")
	res, err := c.instanceService.GetCreateInstanceForm(tenantId, parentExternalId)
	if err != nil {
		log.Println("error getting create instance form: ", err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": res})
}
