package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/adalbertjnr/ws-person/aggregator/client"
	"github.com/adalbertjnr/ws-person/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsume struct {
	consume            *kafka.Consumer
	r                  Replacer
	aggregatorEndpoint client.ClientPicker
}

func NewKafkaConsume(topic string, replacer Replacer, endpoint client.ClientPicker) (*KafkaConsume, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("error initializing kafka consumer client %w", err)
	}

	if err := c.SubscribeTopics([]string{topic}, nil); err != nil {
		log.Fatal(err)
	}
	return &KafkaConsume{
		consume:            c,
		r:                  replacer,
		aggregatorEndpoint: endpoint,
	}, nil
}

func (k *KafkaConsume) start() {
	if _, err := k.consumeDataFromKafka(); err != nil {
		logrus.Error("error", err)
	}
}

func (k *KafkaConsume) consumeDataFromKafka() (*types.Person, error) {
	for {
		msg, err := k.consume.ReadMessage(time.Second)
		if err == nil {
			var data types.Person
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				return nil, fmt.Errorf("error deserializing kafka message from the producer %w", err)
			}
			data.Stage = currentStage
			newDataWithStage := k.r.ReplaceData(data)

			req := &types.AggregatePerson{
				Id:    int32(newDataWithStage.Id),
				Name:  newDataWithStage.Name,
				Age:   int32(newDataWithStage.Age),
				Role:  newDataWithStage.Role,
				Stage: newDataWithStage.Stage,
			}

			if err := k.aggregatorEndpoint.Aggregate(context.TODO(), req); err != nil {
				return nil, err
			}
			//if the error is not about timeout then return it
		} else if !err.(kafka.Error).IsTimeout() {
			return nil, fmt.Errorf("error reading kafka message %w", err)
		}
	}
}
