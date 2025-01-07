package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mcolomerc/confluent_cost_exporter/config"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestGetCosts(t *testing.T) {
	// Mock server to simulate API response
	mockResponse := `{
		"data": [
			{
				"amount": 100.0,
				"discount_amount": 10.0,
				"end_date": "2024-12-16",
				"granularity": "daily",
				"line_type": "usage",
				"network_access_type": "public",
				"original_amount": 110.0,
				"price": 1.0,
				"product": "Kafka",
				"quantity": 100.0,
				"start_date": "2024-12-15",
				"unit": "GB",
				"resource": {
					"display_name": "Test Resource",
					"environment": {
						"id": "env-123"
					},
					"id": "res-123"
				}
			}
		]
	}`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer mockServer.Close()

	// Mock config
	mockConfig := &config.Config{
		Endpoints: config.Endpoints{
			CostsUrl: mockServer.URL,
		},
	}

	// Mock HTTP client
	mockHttpClient := &client.HttpClient{
		Client: mockServer.Client(),
	}

	// Create CostService
	costService := NewCostService(mockHttpClient, mockConfig)

	// Call GetCosts
	costs, err := costService.GetCosts()

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, costs)
	assert.Len(t, costs, 1)

	cost := costs[0]
	assert.Equal(t, 100.0, float64(cost.Amount))
	assert.Equal(t, 10.0, float64(cost.DiscountAmount))
	assert.Equal(t, "2024-12-16", cost.EndDate)
	assert.Equal(t, "daily", cost.Granularity)
	assert.Equal(t, "usage", cost.LineType)
	assert.Equal(t, "public", cost.NetworkAccessType)
	assert.Equal(t, 110.0, float64(cost.OriginalAmount))
	assert.Equal(t, 1.0, float64(cost.Price))
	assert.Equal(t, "Kafka", cost.Product)
	assert.Equal(t, 100.0, float64(cost.Quantity))
	assert.Equal(t, "2024-12-15", cost.StartDate)
	assert.Equal(t, "GB", cost.Unit)
	assert.Equal(t, "Test Resource", cost.Resource.DisplayName)
	assert.Equal(t, "env-123", cost.Resource.Environment.ID)
	assert.Equal(t, "res-123", cost.Resource.ID)
}
