// Package v1 implements routing paths. Each services in own file.
package router

import (
	"time" 
 
	"github.com/gin-gonic/gin"
	"github.com/mcolomerc/confluent_cost_exporter/config"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/controller"
)

type Router struct {
	ExportController controller.ExportController
}

// Routes is the function that defines the routes for the application
func (rt *Router) SetupRouter(cacheExpiration time.Duration) *gin.Engine {
	router := gin.Default()

	// Environment routes
	router.GET("/probe", rt.ExportController.Export(config.PROMETHEUS))
	router.GET("/json", rt.ExportController.Export(config.JSON))

	return router
}
