package routes

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/habib-web-go/gateway-server/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func addAuthRoutes(rg *gin.Engine) (*grpc.ClientConn, error) {
	address := os.Getenv("AUTH_ADDR")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	auth := rg.Group("/auth")
	client := pb.NewAuthServiceClient(conn)
	auth.GET("/req_pq", func(c *gin.Context) {
		req, err := parseReqPQParams(c)
		if err != nil {
			return
		}
		res, err := client.ReqPQ(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, parseReqPQResponse(res))
	})

	auth.GET("/req_dh_params", func(c *gin.Context) {
		req, err := parseReqDHParams(c)
		if err != nil {
			return
		}
		res, err := client.ReqDHParams(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, parseReqDHResponse(res))
	})

	return conn, err
}

func parseReqPQParams(c *gin.Context) (*pb.ReqPQRequest, error) {
	nonce := c.Param("nonce")
	messageId, err := strconv.ParseUint(c.Param("message_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return nil, err
	}
	req := &pb.ReqPQRequest{Nonce: nonce, MessageId: messageId}
	return req, nil
}

func parseReqDHParams(c *gin.Context) (*pb.ReqDHParamsRequest, error) {
	nonce := c.Param("nonce")
	serverNonce := c.Param("server_nonce")
	a := c.Param("a")
	messageId, err := strconv.ParseUint(c.Param("message_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return nil, err
	}
	req := &pb.ReqDHParamsRequest{
		Nonce:       nonce,
		ServerNonce: serverNonce,
		A:           a,
		MessageId:   messageId,
	}
	return req, nil
}

func parseReqPQResponse(res *pb.ReqPQResponse) gin.H {
	return gin.H{
		"g":            fmt.Sprint(res.G),
		"p":            fmt.Sprint(res.P),
		"message_id":   fmt.Sprint(res.MessageId),
		"nonce":        fmt.Sprint(res.Nonce),
		"server_nonce": fmt.Sprint(res.ServerNonce),
	}
}

func parseReqDHResponse(res *pb.ReqDHParamsResponse) gin.H {
	return gin.H{
		"nonce":        fmt.Sprint(res.Nonce),
		"server_nonce": fmt.Sprint(res.ServerNonce),
		"message_id":   fmt.Sprint(res.MessageId),
		"b":            fmt.Sprint(res.B),
	}
}
