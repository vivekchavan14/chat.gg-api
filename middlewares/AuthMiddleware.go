package middlewares

import (
	"net/http"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing token",
			})
			return
		}

		claims, err := utils.ParseJWT(tokenString, jwtKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}
