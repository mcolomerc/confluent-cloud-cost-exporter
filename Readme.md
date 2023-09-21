# Confluent Cloud Metrics and Costs - Prometheus + Grafana

* [Prometheus](https://prometheus.io/) is an open-source systems monitoring and alerting toolkit
* [Grafana](https://grafana.com/) is an open-source platform for monitoring and observability
* [Confluent Cloud Metrics API](https://docs.confluent.io/cloud/current/metrics-api.html) provides metrics for Confluent Cloud resources
* [Confluent Cloud Costs API](https://docs.confluent.io/cloud/current/billing/overview.html) provides costs for Confluent Cloud resources

## Build params with prometheus target labels

* Cost data can take up to 72 hours to become available
* Start date can reach a maximum of one year into the past
* One month is the maximum window between start and end dates
* Period - Current Month  

From: [Confluent Costs API](https://docs.confluent.io/cloud/current/billing/overview.html#retrieve-costs-for-a-range-of-dates)

## Prometheus

Prometheus interval is set to `5m` and timeout to `30s`. More than a `5m` interval is not recommended for Prometheus.

```yaml
- job_name: json_exporter
    scrape_interval: 5m
    scrape_timeout: 30s
    honor_labels: true
    ... 
```

### Scrape Targets

* Confluent Cloud Metrics API

Based on [JSON Exporter](https://github.com/prometheus-community/json_exporter)

```yaml
scrape_configs:
  - job_name: Confluent Cloud
    scrape_interval: 1m
    scrape_timeout: 1m
    honor_timestamps: true
    static_configs:
      - targets:
        - api.telemetry.confluent.cloud
    scheme: https
    basic_auth:
      username: <CONFLUENT_CLOUD_API_KEY>
      password: <CONFLUENT_CLOUD_API_SECRET>
    metrics_path: /v2/metrics/cloud/export
    params:
      resource.kafka.id:
        - <CLUSTER_ID>
```

* Confluent Cloud Costs API

```yaml
- job_name: json_exporter
    scrape_interval: 5m
    scrape_timeout: 30s
    honor_labels: true 
    metrics_path: /probe  
    static_configs:
      - targets: 
        - https://api.confluent.cloud/billing/v1/costs
           
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target 
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: json_exporter:7979     
```

## JSON Prometheus Exporter

Based on [JSON Exporter](https://github.com/prometheus-community/json_exporter). 

Patches:

* https://github.com/prometheus-community/json_exporter/issues/148
  
* Request parameters `start_date` and `end_date` are calculated based on the current month

JSON Exporter builds the Request to Confluent Cloud API and builds the response in Prometheus format.

* Period - Current Month (calculates current month, start and end dates, to build the request parameters `start_date` and `end_date`)
  
* Manages the authentication with Confluent Cloud API, it requires the following environment variables:
  
  * user=`CCLOUD_API_KEY`
  * pass=`CCLOUD_API_SECRET`

Metric definition: `.prom-json-exporter/config.yml`

```yaml
 metrics:
      - name: confluent_cloud_cost
        type: object
        help: Confluent Cloud Resource costs
        path: '{.data[*]}'
        labels:
          id: '{.resource.id}'
          resource: '{.resource.display_name}'
          environment: '{.resource.environment.id}'
          unit: '{.unit}'
          product: '{.product}'
          start: '{.start_date}' 
          end_date: '{.end_date}' 
          granularity: '{.granularity}'  
          discount: '{.discount_amount}'
          price: '{.price}'
          original_amount: '{.original_amount}'
          quantity: '{.quantity}'
        values:
          amount: '{.amount}'  
``` 

* Prometheus query example:

`sum(confluent_cloud_cost_amount)`or `sum(confluent_cloud_cost_amount{id=~"lksqlc.*"})`

### Build the JSON Prometheus Exporter
  
```bash
  docker-compose build
```

## Run

Define required environment variables at `docker-compose.yaml` file.

```bash
  docker-compose up -d
```

### Grafana

Open `localhost:3000` and login with `admin:admin`

Provisioned dashboards:

* Confluent Cloud metrics
* Confluent Cloud costs

![alt text](./docs/Grafana.png) 

### Prometheus

Open `localhost:9090` and check the targets.

## TODOs

* [X] JSON Exporter. Map Confluent Cloud Cloud Cost API response to Prometheus format
* [X] JSON Exporter. Add authentication to Confluent Cloud API
* [X] JSON Exporter. Add start_date and end_date parameters to Confluent Cloud API request
* [ ] JSON Exporter. Reduce the number of request to Confluent Cloud API to get the Costs data.  
* [ ] Grafana. Add more Panels to the dashboards
* [ ] Combine Metrics and Costs in the same dashboard  
* [ ] Alerting