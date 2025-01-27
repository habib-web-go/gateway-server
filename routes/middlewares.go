package routes

import (
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/habib-web-go/gateway-server/gen/grpc/auth"
	redis "github.com/habib-web-go/gateway-server/redis"
)

type authCheckStruct struct {
	AuthKey string `json:"authKey,omitempty"`
}

func rateLimit(ctx *gin.Context) {
	ip := ctx.ClientIP()
	rate, err := redis.GetRateLimit(ip)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
		return
	}

	if rate <= 0 {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many request. Come back later.",
		})
		ctx.Abort()
		return
	}

	redis.DecreaseRateLimit(ip)
	ctx.Next()
}

func createAuthCheckMiddleware(client pb.AuthServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var authkeyStruct *authCheckStruct
		err := ctx.ShouldBindBodyWith(&authkeyStruct, binding.JSON)
		if err != nil {
			ctx.JSON(http.StatusNetworkAuthenticationRequired, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}
		authkey := authkeyStruct.AuthKey
		req := pb.IsValidAuthKeyRequest{Authkey: authkey}
		res, err := client.IsValidAuthkey(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}
		if !res.IsValid {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid authkey",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}

}
