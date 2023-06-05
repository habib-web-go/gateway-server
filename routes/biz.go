package routes

import (
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	authpb "github.com/habib-web-go/gateway-server/gen/grpc/auth"
	pb "github.com/habib-web-go/gateway-server/gen/grpc/biz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func addBizRoutes(
	rg *gin.Engine,
	authClient authpb.AuthServiceClient,
) (*grpc.ClientConn, pb.SQLServiceClient, error) {
	address := os.Getenv("BIZ_ADDR")
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	biz := rg.Group("/biz")

	client := pb.NewSQLServiceClient(conn)
	biz.Use(createAuthCheckMiddleware(authClient))

	biz.POST("/get_users", getUsersHandler(client))
	biz.POST("/get_users_with_sql_inject", getUsersWithSqlInjectHandler(client))

	return conn, client, nil
}

// @Summary Get users
// @Description Retrieves a list of users from the database
// @Accept json
// @Produce json
// @Param req body pb.GetUsersRequest true "Request body"
// @Success 200 {object} pb.GetUsersResponse
// @Failure 400
// @Failure 500
// @Router /biz/get_users [post]
func getUsersHandler(client pb.SQLServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req *pb.GetUsersRequest
		if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
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
	}
}

// @Summary Get users with SQL injection vulnerability
// @Description Retrieves a list of users from the database using a SQL injection vulnerable query
// @Accept json
// @Produce json
// @Param req body pb.GetUsersWithSqlInjectRequest true "Request body"
// @Success 200 {object} pb.GetUsersResponse
// @Failure 400
// @Failure 500
// @Router /biz/get_users_with_sql_inject [post]
func getUsersWithSqlInjectHandler(client pb.SQLServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req *pb.GetUsersWithSqlInjectRequest
		if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
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
	}
}
