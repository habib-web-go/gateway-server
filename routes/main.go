package routes

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var router = gin.Default()

func Run() {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authConn, authClient, err := addAuthRoutes(router)
	defer authConn.Close()
	if err != nil {
		log.Fatalf("Auth did not connect: %v", err)
		return
	}

	bizConn, _, err := addBizRoutes(router, authClient)
	defer bizConn.Close()
	if err != nil {
		log.Fatalf("Biz did not connect: %v", err)
		return
	}

	port := os.Getenv("GIN_PORT")
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
