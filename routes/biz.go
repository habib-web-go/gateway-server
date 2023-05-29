package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func addBizRoutes(rg *gin.Engine) (*grpc.ClientConn, error) {
	address := os.Getenv("BIZ_ADDR")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	biz := rg.Group("/auth")

	biz.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
	return conn, err
}
