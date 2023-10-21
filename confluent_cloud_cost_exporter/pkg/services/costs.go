package services

import (
	"log"
	"strconv"

	"net/url"
	"time"

	"github.com/mcolomerc/confluent_cost_exporter/config"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/client"
	"github.com/mitchellh/mapstructure"
)

type Cost struct { // This is the struct that will hold the data from the API
	Amount            KeepZero `json:"amount" mapstructure:"amount"`
	DiscountAmount    KeepZero `json:"discount_amount" mapstructure:"discount_amount"`
	EndDate           string   `json:"end_date" mapstructure:"end_date"`
	Granularity       string   `json:"granularity" mapstructure:"granularity"`
	LineType          string   `json:"line_type" mapstructure:"line_type"`
	NetworkAccessType string   `json:"network_access_type" mapstructure:"network_access_type"`
	OriginalAmount    KeepZero `json:"original_amount" mapstructure:"original_amount"`
	Price             KeepZero `json:"price" mapstructure:"price"`
	Product           string   `json:"product" mapstructure:"product"`
	Quantity          KeepZero `json:"quantity" mapstructure:"quantity"`
	StartDate         string   `json:"start_date" mapstructure:"start_date"`
	Unit              string   `json:"unit" mapstructure:"unit"`
	Resource          `json:"resource"  mapstructure:"resource"`
}

type Resource struct {
	DisplayName string `json:"display_name" mapstructure:"display_name"`
	Environment `json:"environment" mapstructure:"environment"`
	ID          string `json:"id" mapstructure:"id" `
}

type Environment struct {
	ID string `json:"id" mapstructure:"id"`
}

type CostService struct {
	Client *client.HttpClient
	Config *config.Config
}

type KeepZero float64

func (f KeepZero) MarshalJSON() ([]byte, error) {
	if float64(f) == float64(int(f)) {
		return []byte(strconv.FormatFloat(float64(f), 'f', 1, 32)), nil
	}
	return []byte(strconv.FormatFloat(float64(f), 'f', -1, 32)), nil
}

func NewCostService(client *client.HttpClient, config *config.Config) *CostService {
	return &CostService{
		Client: client,
		Config: config,
	}
}
func (e *CostService) GetCosts() ([]Cost, error) {
	endpoint := e.Config.Endpoints.CostsUrl
	baseURL, _ := url.Parse(endpoint)
	params := url.Values{}
	// Cost data can take up to 72 hours to become available
	// Start date can reach a maximum of one year into the past
	// One month is the maximum window between start and end dates
	// Period - Current month
	_, currentMonth, _ := time.Now().Date()
	startDate, endDate := getDates(int(currentMonth))
	params.Add("end_date", endDate)
	params.Add("start_date", startDate)
	baseURL.RawQuery = params.Encode()
	response, err := e.Client.GetData(baseURL.String())
	if err != nil {
		log.Println("Error:", err)
		return nil, nil
	}
	var costs []Cost
	for _, env := range response.Data {
		var result Cost
		err := mapstructure.Decode(env, &result)
		if err != nil {
			log.Println("Error:", err)
		}

		costs = append(costs, result)
	}
	defer timer("Build Cost from Confluent Cloud API")()
	return costs, nil
}

func getDates(month int) (string, string) {
	currentMonth := time.Month(month)
	now := time.Now()
	currentYear, _, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return firstOfMonth.Format("2006-01-02"), lastOfMonth.Format("2006-01-02")
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", name, time.Since(start))
	}
}
