package router

import (
	delivery "github.com/sigit14ap/shop-service/internal/delivery/http"
	"github.com/sigit14ap/shop-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(shopHandler *delivery.ShopHandler) *gin.Engine {
	router := gin.New()

	v1 := router.Group("/api/v1/shop").Use(middleware.ServiceMiddleware())
	{
		v1.POST("/register", shopHandler.Register)
		v1.POST("/login", shopHandler.Login)

		shop := v1.Use(middleware.AuthMiddleware())
		{
			shop.GET("/me", shopHandler.Me)
		}
	}

	return router
}
