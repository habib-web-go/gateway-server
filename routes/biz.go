package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	pb "github.com/habib-web-go/gateway-server/grpc/biz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func addBizRoutes(rg *gin.Engine) (*grpc.ClientConn, *gin.RouterGroup, error) {
	address := os.Getenv("BIZ_ADDR")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	biz := rg.Group("/auth")
	client := pb.NewSQLServiceClient(conn)
	biz.GET("/get_users", func(c *gin.Context) {
		var req *pb.GetUsersRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.GetUsers(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	biz.GET("/get_users_with_sql_inject", func(c *gin.Context) {
		var req *pb.GetUsersWithSqlInjectRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		res, err := client.GetUsersWithSqlInject(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
	})
	return conn, biz, nil
}
