package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	pb "github.com/habib-web-go/gateway-server/gen/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func addAuthRoutes(rg *gin.Engine) (*grpc.ClientConn, pb.AuthServiceClient, error) {
	// Get the address of the authentication service from environment variable
	address := os.Getenv("AUTH_ADDR")

	// Establish a gRPC connection to the authentication service
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	// Create a new group for authentication routes
	auth := rg.Group("/auth")

	// Create a rate limiter middleware for the authentication routes
	auth.Use(rateLimit)

	// Create a client for the authentication service using the gRPC connection
	client := pb.NewAuthServiceClient(conn)

	// Add a route for the ReqPQ endpoint
	auth.POST("/req_pq", handleReqPQ(client))

	// Add a route for the ReqDHParams endpoint
	auth.POST("/req_dh_params", handleReqDHParams(client))

	return conn, client, nil
}

// @Summary Request PQ
// @Description Request PQ from server
// @Tags Auth
// @Accept json
// @Produce json
// @Param req body pb.ReqPQRequest true "Request Params"
// @Success 200 {object} pb.ReqPQResponse
// @Failure 400
// @Router /auth/req_pq [post]
func handleReqPQ(client pb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}

// @Summary Request DH Params
// @Description Request Diffie-Hellman parameters from server
// @Tags Auth
// @Accept json
// @Produce json
// @Param req body pb.ReqDHParamsRequest true "Request Params"
// @Success 200 {object} pb.ReqDHParamsResponse
// @Failure 400
// @Router /auth/req_dh_params [post]
func handleReqDHParams(client pb.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}
