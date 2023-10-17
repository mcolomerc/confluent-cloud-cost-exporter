package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mcolomerc/confluent_cost_exporter/config"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/exporters"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/services"
)

type ExportController struct {
	CostsService *services.CostService
	Exporters    map[config.Format]exporters.Exporter
}

func NewExportController(costsService *services.CostService, exporters map[config.Format]exporters.Exporter) *ExportController {
	return &ExportController{
		CostsService: costsService,
		Exporters:    exporters,
	}
}

func (controller *ExportController) Export(format config.Format) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer timer("Export")()
		response, err := controller.CostsService.GetCosts()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error getting environments"})
			return
		}
		controller.Exporters[format].ExportCosts(response, c)
		return
	}
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", name, time.Since(start))
	}
}
