// Package app configures and runs application.
package app

import (
	"log"

	"github.com/mcolomerc/confluent_cost_exporter/config"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/client"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/controller"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/exporters"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/router"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/services"
)

type App struct {
	Config *config.Config
}

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	// http client
	client := client.NewHttpClient(cfg.Credentials)
	// controller service
	costsService := services.NewCostService(client, cfg)

	var allExporters map[config.Format]exporters.Exporter = make(map[config.Format]exporters.Exporter)
	// exporters
	allExporters[config.PROMETHEUS] = exporters.NewPromExporter(cfg.PromConfig)
	allExporters[config.JSON] = exporters.NewJSONExporter(*cfg)
	// controller
	controller := controller.NewExportController(costsService, allExporters)
	// router
	router := router.Router{
		ExportController: *controller,
	}
	log.Printf("Cache expiration : %v", cfg.Cache.Expiration)
	//setup routes
	r := router.SetupRouter(cfg.Cache.Expiration)

	// running
	r.Run()

}
