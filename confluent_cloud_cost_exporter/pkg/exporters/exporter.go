package exporters

import (
	"github.com/gin-gonic/gin"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/services"
)

type Exporter interface {
	ExportCosts(costs []services.Cost, c *gin.Context) error
	OutputType() OutputType
}

// Create a enum for the different types of exporters
type OutputType int

const (
	TEXT OutputType = iota
	JSON
)
