package exporters

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/mcolomerc/confluent_cost_exporter/config"
	model "github.com/mcolomerc/confluent_cost_exporter/pkg/avro"
	"github.com/mcolomerc/confluent_cost_exporter/pkg/services"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
)

type KafkaExporter struct {
	CostsService *services.CostService
	Config       *config.Config
	Producer     *kafka.Producer
	Serializer   *avro.SpecificSerializer
}

func NewKafkaExporter(costsService *services.CostService, cfg config.Config) (*KafkaExporter, error) {
	log.Println("ProduceCosts")

	kafkaConfig := cfg.Cron.Target.KafkaConfig
	// Produce a new record to the topic...
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.Bootstrap,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_SSL",
		"sasl.username":     kafkaConfig.Credentials.Key,
		"sasl.password":     kafkaConfig.Credentials.Secret,
		"batch.size":        64000,
		"linger.ms":         1000,
		"acks":              "all",
	})
	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		return nil, err
	}

	client, err := schemaregistry.NewClient(schemaregistry.NewConfigWithAuthentication(
		kafkaConfig.SchemaRegistry.Endpoint,
		kafkaConfig.SchemaRegistry.SR_Credentials.Key,
		kafkaConfig.SchemaRegistry.SR_Credentials.Secret))

	if err != nil {
		log.Printf("Failed to create schema registry client: %s\n", err)
		return nil, err
	}

	ser, err := avro.NewSpecificSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())
	if err != nil {
		log.Printf("Failed to create serializer: %s\n", err)
		return nil, err
	}
	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to Topic[partition]@offset: %v\n", ev.TopicPartition)
				}
			}
		}
		producer.Flush(2 * 1000) // 2 seconds
	}()

	return &KafkaExporter{
		Config:       &cfg,
		CostsService: costsService,
		Producer:     producer,
		Serializer:   ser,
	}, nil
}

func (exp KafkaExporter) ExportCosts(done chan bool) error {
	// Schedule jobs
	// Get cost every hour
	costsJob := gocron.NewScheduler(time.UTC)
	costsJob.Cron(exp.Config.Cron.Expression).Do(func() {
		log.Println("Starting costs job...")
		response, err := exp.CostsService.GetCosts()
		if err != nil {
			log.Printf("Error %s", "error getting costs")
			return
		}
		exp.ProduceCosts(response)

	})
	costsJob.StartBlocking() // blocks current goroutine
	return nil
}

func (exp KafkaExporter) OutputType() OutputType {
	return JSON
}

func (exp KafkaExporter) ProduceCosts(costs []services.Cost) error {
	topic := exp.Config.Cron.Target.KafkaConfig.Topic
	log.Printf("Producing %d messages to topic %s\n", len(costs), topic)

	for _, cost := range costs {
		avroCost, _ := transformCosts(cost)
		payload, err := exp.Serializer.Serialize(topic, &avroCost)
		if err != nil {
			log.Printf("Failed to serialize message: %s\n", err)
			return err
		}
		exp.Producer.Produce(
			&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &topic,
					Partition: kafka.PartitionAny},
				Key:   []byte(cost.EndDate),
				Value: payload,
			}, nil)
	}
	defer timer("Produce to Kafka")()
	return nil
}

func transformCosts(in services.Cost) (model.Cost, error) {
	var out model.Cost = model.Cost{}
	out.Amount = float64(in.Amount)
	out.Discount_amount = float64(in.DiscountAmount)
	out.End_date = in.EndDate
	out.Granularity = in.Granularity
	out.Line_type = in.LineType
	out.Network_access_type = in.NetworkAccessType
	out.Original_amount = float64(in.OriginalAmount)
	out.Price = float64(in.Price)
	out.Product = in.Product
	out.Quantity = float64(in.Quantity)
	out.Start_date = in.StartDate
	out.Unit = in.Unit
	var resource model.Resource = model.Resource{}
	resource.Display_name = in.Resource.DisplayName
	var environment model.Environment = model.Environment{}
	environment.Id = in.Resource.Environment.ID
	resource.Environment = environment
	resource.Id = in.Resource.ID
	out.Resource = resource
	return out, nil
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %v\n", name, time.Since(start))
	}
}
