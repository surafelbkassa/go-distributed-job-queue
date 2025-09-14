package routers

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/surafelbkassa/go-distributed-job-queue/Delivery/Controllers"
)

func Init(gin *gin.Engine) {
	Routes := gin.Group("")
	// Register endpoint
	Routes.POST("/auth/register", controllers.RegisterHandler)
}
