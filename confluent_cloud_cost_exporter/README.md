# Confluent Cloud Cost Exporter

Exporters:

* `/probe`: Confluent Cloud Cost as [prometheus](https://prometheus.io/)
* `/json`: Confluent Cloud Cost as [json](https://www.json.org/json-en.html)  

## Configuration

The exporter is configured via a YAML file, by default located at `./config.yml`. 

The following configuration options are available:

```yaml
credentials: 
  key: <CONFLUENT_API_KEY>
  secret: <CONFLUENT_API_SECRET>
cache:
  expiration: 1200m

```

or as environment variables:

* CONFLUENT_CLOUD_API_KEY=<CONFLUENT_CLOUD_API_KEY>
* CONFLUENT_CLOUD_API_SECRET=<CONFLUENT_CLOUD_API_SECRET>
* CACHE_EXPIRATION=240m #duration

## Exporters

### JSON Exporter

Sample result:

```json
[
  {
    "amount": 0.9576,
    "discount_amount": 0.0,
    "end_date": "2023-10-02",
    "granularity": "DAILY",
    "line_type": "CONNECT_NUM_TASKS",
    "network_access_type": "INTERNET",
    "original_amount": 0.9576,
    "price": 0.03993055,
    "product": "CONNECT",
    "quantity": 24.0,
    "start_date": "2023-10-01",
    "unit": "Task-hour",
    "resource": {
      "display_name": "DatagenSource_1",
      "environment": {
        "id": "env-xxxxxx"
      },
      "id": "lcc-xxxxxx"
    }, 
  } 
]
```

### Prometheus exporter

Sample result:

```txt
 # HELP confluent_cloud_cost_amount Confluent Cloud Resource costs
# TYPE confluent_cloud_cost_amount untyped
confluent_cloud_cost_amount{discount="0",end_date="2023-10-02",environment="env-xxxxxx",granularity="DAILY",id="lkc-xxxxxx",original_amount="0",price="0.0001326",product="KAFKA",quantity="2.6151538e-05",resource="connect",start="2023-10-01",unit="GB-hour"} 0
confluent_cloud_cost_amount{discount="0",end_date="2023-10-02",environment="env-xxxxxx",granularity="DAILY",id="lkc-xxxxxx",original_amount="3.1368",price="0.00484",product="KAFKA",quantity="648",resource="connect",start="2023-10-01",unit="Partition-hour"} 3.1368
```

## Development

```sh
go run main.go 
```

## TODO

* [ ] Add Kafka topic exporter. 