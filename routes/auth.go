package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	pb "github.com/habib-web-go/gateway-server/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func addAuthRoutes(rg *gin.Engine) (*grpc.ClientConn, *gin.RouterGroup, error) {
	address := os.Getenv("AUTH_ADDR")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	auth := rg.Group("/auth")
	client := pb.NewAuthServiceClient(conn)
	auth.GET("/req_pq", func(c *gin.Context) {
		var req *pb.ReqPQRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.ReqPQ(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	auth.GET("/req_dh_params", func(c *gin.Context) {
		var req *pb.ReqDHParamsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.ReqDHParams(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	return conn, auth, nil
}
