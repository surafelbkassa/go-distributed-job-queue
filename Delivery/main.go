package main

import (
	"github.com/gin-gonic/gin"
	routers "github.com/surafelbkassa/go-distributed-job-queue/Delivery/Routers"
)

func main() {
	router := gin.Default()
	routers.Init(router)
	router.Run()
}
