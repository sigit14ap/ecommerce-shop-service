package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sigit14ap/shop-service/helpers"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		if token == "" {
			helpers.ErrorResponse(context, http.StatusUnauthorized, "Unauthorized")
			context.Abort()
			return
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			context.Abort()
			return
		}

		tokenString := parts[1]

		tokenData, err := helpers.ParseJWT(tokenString)

		if err != nil {
			helpers.ErrorResponse(context, http.StatusUnauthorized, "Unauthorized")
			context.Abort()
			return
		}

		context.Set("shopID", tokenData.ShopID)
		context.Next()
	}
}
