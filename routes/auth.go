package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func addAuthRoutes(rg *gin.Engine) (*grpc.ClientConn, error) {
	address := os.Getenv("AUTH_ADDR")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	auth := rg.Group("/auth")

	auth.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	return conn, err
}
