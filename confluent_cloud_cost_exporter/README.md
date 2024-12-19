# Confluent Cloud Cost Exporter

## Configuration

The exporter is configured via a YAML file, by default located at `./config.yml`.

`confluent_cloud_cost_exporter -config=<path/to/config.yml>`

The following configuration options are available:

```yaml
credentials: 
  key: <CONFLUENT_API_KEY>
  secret: <CONFLUENT_API_SECRET> 
```

or as environment variables:

* CONFLUENT_CLOUD_API_KEY=<CONFLUENT_CLOUD_API_KEY>
* CONFLUENT_CLOUD_API_SECRET=<CONFLUENT_CLOUD_API_SECRET>

## Exporters

### Web Exporters

Default exporter. Web exporter will expose cost information in the following endpoints:

* `/probe`: Confluent Cloud Cost as [prometheus](https://prometheus.io/)
* `/json`: Confluent Cloud Cost as [json](https://www.json.org/json-en.html)  

Configure cache expiration. Default: `30m`

```yaml
web:
  cache:
    expiration: 1200m
```

or setting the environment variable `CACHE_EXPIRATION`=<CACHE_EXPIRATION>

#### JSON Exporter

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

#### Prometheus exporter

Sample result:

```txt
 # HELP confluent_cloud_cost_amount Confluent Cloud Resource costs
# TYPE confluent_cloud_cost_amount untyped
confluent_cloud_cost_amount{discount="0",end_date="2023-10-02",environment="env-xxxxxx",granularity="DAILY",id="lkc-xxxxxx",original_amount="0",price="0.0001326",product="KAFKA",quantity="2.6151538e-05",resource="connect",start="2023-10-01",unit="GB-hour"} 0
confluent_cloud_cost_amount{discount="0",end_date="2023-10-02",environment="env-xxxxxx",granularity="DAILY",id="lkc-xxxxxx",original_amount="3.1368",price="0.00484",product="KAFKA",quantity="648",resource="connect",start="2023-10-01",unit="Partition-hour"} 3.1368
```

### Cron Exporters

Cron Job that gets Confluent Cloud billing information and exports it to a given target.

If cron is configured, the exporter will run as a cron job only. If no cron is configured, the exporter will run as a web server.

```yaml
cron:
  expression: "*/5 * * * *" 
```

or setting the environment variable `CRON_EXPRESSION`=<CRON_EXPRESSION>  

#### Target: Confluent Cloud Topic

Act as a cron job that gets Confluent Cloud billing information and exports it to a given topic in Confluent Cloud.
It defines and publish to the Confluent Cloud Schema Registry an AVRO schema to serialize the data.
Destination Topic should be created on Confluent Cloud before running the exporter.

```yaml
credentials: 
  key: <CONFLUENT_CLOUD_API_KEY>
  secret: <CONFLUENT_CLOUD_API_SECRET>
cron:
  expression: "<CRON_EXPRESSION>" 
  target:
    kafka:
      bootstrap: <BOOTSTRAP_SERVER>
      credentials: 
        key: <CONFLUENT_CLOUD_KAFKA_API_KEY>
        secret: <CONFLUENT_CLOUD_KAFKA_API_SECRET>
      schemaRegistry:
        endpoint: <SCHEMA_REGISTRY_ENDPOINT>
        credentials:
          key: <SCHEMA_REGISTRY_API_KEY>
          secret: <SCHEMA_REGISTRY_API_SECRET>
      topic: <TOPIC_NAME>
```

or setting environment variables:

* **BOOTSTRAP**=<CONFLUENT_CLOUD_BOOTSTRAP_SERVER>
* **KAFKA_API_KEY**=<CONFLUENT_KAFKA_CLOUD_API_KEY>
* **KAFKA_API_SECRET**=<CONFLUENT_KAFKA_CLOUD_API_SECRET>
* **TOPIC**=<TOPIC_NAME>
* **SCHEMA_REGISTRY_ENDPOINT**=<SCHEMA_REGISTRY_ENDPOINT>
* **SCHEMA_REGISTRY_API_KEY**=<SCHEMA_REGISTRY_API_KEY>
* **SCHEMA_REGISTRY_API_SECRET**=<SCHEMA_REGISTRY_API_SECRET>

## Development

### Run locally

With the default configuration:

```sh
go run main.go 
```

Providing a configuration file:

```sh
go run main.go -config=<path/to/config.yml>
```

### Docker

Build:

```sh
docker build -t mcolomerc/confluent_cloud_cost_exporter .
```

**Run**

Using environment variables:

Web exporters:

```sh
docker run -p 7979:7979 \ 
--env CONFLUENT_CLOUD_API_KEY='<CONFLUENT_CLOUD_API_KEY>' \
--env CONFLUENT_CLOUD_API_SECRET='<CONFLUENT_CLOUD_API_SECRET>' \
mcolomerc/confluent_cloud_cost_exporter 
```

Using a configuration file:

```sh
docker run -p 7979:7979 -v <path/to/config.yml>:/bin/config.yml confluent/confluent_cloud_cost_exporter 
```

Test endpoint:

```sh
curl http://localhost:7979/json
```

### AVRO

`go generate github.com/mcolomerc/confluent_cost_exporter/generate`
