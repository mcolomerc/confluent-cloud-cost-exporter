package services

import (
	"log"
	"strconv"

	"net/url"
	"time"

	"github.com/mcolomerc/confluent_cost_exporter/config"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/cache"
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
	Client   *client.HttpClient
	Config   *config.Config
	CacheTTL *cache.TTLCache[string, []Cost]
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
		Client:   client,
		Config:   config,
		CacheTTL: cache.NewTTL[string, []Cost](),
	}
}

/*
Cost data can take up to 72 hours to become available
Start date can reach a maximum of one year into the past
One month is the maximum window between start and end dates
Period - Window size: 1 day - 3 days ago from today - if today is 2024-12-19 then the period is 2024-12-15 to 2024-12-16
*/
func (e *CostService) GetCosts() ([]Cost, error) {
	endpoint := e.Config.Endpoints.CostsUrl
	baseURL, _ := url.Parse(endpoint)
	params := url.Values{}

	startDate, endDate := getDatesFromInterval(e.Config.Period.DaysAgo, e.Config.Period.Window)
	params.Add("end_date", endDate)
	params.Add("start_date", startDate)
	baseURL.RawQuery = params.Encode()

	fromCache, found := e.getFromCache(baseURL.String())
	if found {
		log.Println("Get from cache")
		return fromCache, nil
	}
	log.Println("Request data from Confluent Cloud API")
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

	e.CacheTTL.Set(baseURL.String(), costs, e.Config.Web.Expiration)

	defer timer("Build Cost from Confluent Cloud API")()
	return costs, nil
}

func (e *CostService) getFromCache(key string) ([]Cost, bool) {
	costs, found := e.CacheTTL.Get(key)
	if found {
		return costs, true
	}
	return nil, false
}

func getDatesFromInterval(daysAgo int, wSize int) (string, string) {
	now := time.Now()

	currentLocation := now.Location()

	// Calculate the end date (daysAgo days before today)
	endDate := now.AddDate(0, 0, -daysAgo).In(currentLocation)

	// Calculate the start date (daysAgo + interval days before today)
	startDate := endDate.AddDate(0, 0, -wSize).In(currentLocation)

	// Format the dates as strings
	return startDate.Format("2006-01-02"), endDate.Format("2006-01-02")
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", name, time.Since(start))
	}
}
