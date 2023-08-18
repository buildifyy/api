package instance

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Controller interface {
	GetCreateInstanceForm(c *gin.Context)
}

type controller struct {
	instanceService Service
}

func NewController(instanceService Service) Controller {
	return &controller{
		instanceService: instanceService,
	}
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
