package main

import (
	"encoding/json"
	"fmt"

	"github.com/adalbertjnr/ws-person/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer interface {
	ProduceToKafka(data types.Person) error
}

type KafkaProduce struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProduce(topic string) (*KafkaProduce, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, fmt.Errorf("error initializing kafka producer client %w", err)
	}
	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &KafkaProduce{
		producer: p,
		topic:    topic,
	}, nil
}

func (k *KafkaProduce) ProduceToKafka(data types.Person) error {
	d, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling the data before produce to kafka %w", err)
	}
	return k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.topic, Partition: kafka.PartitionAny},
		Value:          []byte(d),
	}, nil)
}
