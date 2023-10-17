package exporters

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mcolomerc/confluent_cost_exporter/config"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/services"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromExporter struct {
	Config *config.PromConfig
}

type JSONProm struct {
	Data []services.Cost `json:"data"`
}

func NewPromExporter(cfgProm config.PromConfig) *PromExporter {
	return &PromExporter{
		Config: &cfgProm,
	}
}

func (service *PromExporter) OutputType() OutputType {
	return TEXT
}

func (service *PromExporter) ExportCosts(costs []services.Cost, c *gin.Context) error {
	registry := prometheus.NewPedanticRegistry()

	metrics, err := service.CreateMetricsList(service.Config.Modules["default"])
	if err != nil {
		log.Println("Failed to create metrics list from config", "err", err)
	}
	jsonMetricCollector := JSONMetricCollector{JSONMetrics: metrics}

	jsonData := JSONProm{Data: costs}

	b, err := json.Marshal(jsonData)
	if err != nil {
		fmt.Println(err)
	}
	jsonMetricCollector.Data = b

	registry.MustRegister(jsonMetricCollector)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(c.Writer, c.Request)
	return nil
}

func (service *PromExporter) CreateMetricsList(c config.Module) ([]JSONMetric, error) {
	var (
		metrics   []JSONMetric
		valueType prometheus.ValueType
	)
	for _, metric := range c.Metrics {
		switch metric.ValueType {
		case config.ValueTypeGauge:
			valueType = prometheus.GaugeValue
		case config.ValueTypeCounter:
			valueType = prometheus.CounterValue
		default:
			valueType = prometheus.UntypedValue
		}
		switch metric.Type {
		case config.ValueScrape:
			var variableLabels, variableLabelsValues []string
			for k, v := range metric.Labels {
				variableLabels = append(variableLabels, k)
				variableLabelsValues = append(variableLabelsValues, v)
			}
			jsonMetric := JSONMetric{
				Type: config.ValueScrape,
				Desc: prometheus.NewDesc(
					metric.Name,
					metric.Help,
					variableLabels,
					nil,
				),
				KeyJSONPath:            metric.Path,
				LabelsJSONPaths:        variableLabelsValues,
				ValueType:              valueType,
				EpochTimestampJSONPath: metric.EpochTimestamp,
			}
			metrics = append(metrics, jsonMetric)
		case config.ObjectScrape:
			for subName, valuePath := range metric.Values {
				name := MakeMetricName(metric.Name, subName)
				var variableLabels, variableLabelsValues []string
				for k, v := range metric.Labels {
					variableLabels = append(variableLabels, k)
					variableLabelsValues = append(variableLabelsValues, v)
				}
				jsonMetric := JSONMetric{
					Type: config.ObjectScrape,
					Desc: prometheus.NewDesc(
						name,
						metric.Help,
						variableLabels,
						nil,
					),
					KeyJSONPath:            metric.Path,
					ValueJSONPath:          valuePath,
					LabelsJSONPaths:        variableLabelsValues,
					ValueType:              valueType,
					EpochTimestampJSONPath: metric.EpochTimestamp,
				}
				metrics = append(metrics, jsonMetric)
			}
		default:
			return nil, fmt.Errorf("Unknown metric type: '%s', for metric: '%s'", metric.Type, metric.Name)
		}
	}
	return metrics, nil
}

func MakeMetricName(parts ...string) string {
	return strings.Join(parts, "_")
}

func SanitizeValue(s string) (float64, error) {
	var err error
	var value float64
	var resultErr string

	if value, err = strconv.ParseFloat(s, 64); err == nil {
		return value, nil
	}
	resultErr = fmt.Sprintf("%s", err)

	if boolValue, err := strconv.ParseBool(s); err == nil {
		if boolValue {
			return 1.0, nil
		}
		return 0.0, nil
	}
	resultErr = resultErr + "; " + fmt.Sprintf("%s", err)

	if s == "<nil>" {
		return math.NaN(), nil
	}
	return value, fmt.Errorf(resultErr)
}

func SanitizeIntValue(s string) (int64, error) {
	var err error
	var value int64
	var resultErr string

	if value, err = strconv.ParseInt(s, 10, 64); err == nil {
		return value, nil
	}
	resultErr = fmt.Sprintf("%s", err)

	return value, fmt.Errorf(resultErr)
}
