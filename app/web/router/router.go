package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type AppRouter struct {
	*gin.Engine
}

func BuildRouter() *AppRouter {

	if os.Getenv("ENVIRONMENT") == "development" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	// appV1Router := router.Group("/v1")

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Endpoint not found"})
	})

	return &AppRouter{router}
}
