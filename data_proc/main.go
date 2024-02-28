package main

import (
	"log"

	"github.com/adalbertjnr/ws-person/aggregator/client"
)

const (
	kafkaTopic         = "wstopic"
	currentStage       = "data_proc stage"
	httpServerEndpoint = "http://localhost:3000"
)

func main() {
	var (
		r   Replacer
		err error
	)
	r = NewDataMiddlewareLogger(NewReplace())
	consumeAndSendToAggregator, err := NewKafkaConsume(kafkaTopic, r, client.NewEndpoint(httpServerEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	consumeAndSendToAggregator.start()
}
