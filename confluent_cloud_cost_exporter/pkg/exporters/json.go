package exporters

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcolomerc/confluent_cost_exporter/config"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/services"
)

type JSONExporter struct {
	Config *config.Config
}

func NewJSONExporter(cfg config.Config) *JSONExporter {
	return &JSONExporter{
		Config: &cfg,
	}
}

func (exp JSONExporter) ExportCosts(costs []services.Cost, c *gin.Context) error {
	c.JSON(http.StatusOK, costs)
	return nil
}

func (exp JSONExporter) OutputType() OutputType {
	return JSON
}
