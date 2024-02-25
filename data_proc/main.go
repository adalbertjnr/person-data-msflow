package main

import (
	"log"
)

const (
	kafkaTopic   = "wstopic"
	currentStage = "data_proc stage"
)

func main() {
	var (
		r   Replacer
		err error
	)
	r = NewDataMiddlewareLogger(NewReplace())
	consumeAndSendToAggregator, err := NewKafkaConsume(kafkaTopic, r)
	if err != nil {
		log.Fatal(err)
	}
	consumeAndSendToAggregator.start()
}
