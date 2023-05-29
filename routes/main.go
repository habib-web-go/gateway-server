package routes

import (
	"log"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run() {

	authConn, err := addAuthRoutes(router)
	defer authConn.Close()
	if err != nil {
		log.Fatalf("Auth did not connect: %v", err)
		return
	}
	addBizRoutes(router)
	if err != nil {
		log.Fatalf("Biz did not connect: %v", err)
		return
	}
	if err := router.Run(":6433"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
