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
	log.Println("Running Confluent Cost exporter enabled")
	// http client
	client := client.NewHttpClient(cfg.Credentials)
	// controller service
	costsService := services.NewCostService(client, cfg)

	// Push exporters (cron jobs)
	log.Println(cfg.Cron.Expression)
	if cfg.Cron.Expression != "" {
		log.Println("Kafka exporter enabled")
		kafkaExporter, err := exporters.NewKafkaExporter(costsService, *cfg)
		if err != nil {
			log.Println("Error creating kafka exporter")
			return
		}
		log.Printf("Kafka exporter created - Cron: %v", cfg.Cron.Expression)
		done := make(chan bool)
		err = kafkaExporter.ExportCosts(done) //CRON Job
		if err != nil {
			log.Printf("Error %s", "error running kafka exporter")
			return
		}
	} else {
		// web exporters
		var webExporters map[config.Exporter]exporters.Exporter = make(map[config.Exporter]exporters.Exporter)
		// exporters
		webExporters[config.PROMETHEUS] = exporters.NewPromExporter(cfg.Web.PromConfig)
		webExporters[config.JSON] = exporters.NewJSONExporter(*cfg)
		// controller
		controller := controller.NewExportController(costsService, webExporters)
		// router
		router := router.Router{
			ExportController: *controller,
		}
		log.Printf("Cache expiration : %v", cfg.Web.Cache.Expiration)
		//setup routes
		r := router.SetupRouter(cfg.Web.Cache.Expiration)

		// running
		r.Run()
	}
}
