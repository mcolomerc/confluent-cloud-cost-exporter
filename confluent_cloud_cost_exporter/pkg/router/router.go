// Package v1 implements routing paths. Each services in own file.
package router

import (
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
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
	store := persistence.NewInMemoryStore(cacheExpiration)
	// Environment routes
	router.GET("/probe", cache.CachePage(store, cacheExpiration, rt.ExportController.Export(config.PROMETHEUS)))
	router.GET("/json", cache.CachePage(store, cacheExpiration, rt.ExportController.Export(config.JSON)))

	return router
}
