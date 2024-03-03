package main

import (
	"log"

	"github.com/adalbertjnr/ws-person/aggregator/client"
)

const (
	kafkaTopic         = "wstopic"
	currentStage       = "data_proc stage"
	httpServerEndpoint = "http://localhost:3001"
	grpcServerEndpoint = ":4001"
)

func main() {
	var (
		r   Replacer
		err error
	)
	r = NewDataMiddlewareLogger(NewReplace())
	// here I can instanciate both clientes (grpc or http)
	// grpc
	// c, err := client.NewGRPCClientEndpoint(grpcServerEndpoint)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// http
	c := client.NewHTTPClientEndpoint(httpServerEndpoint)
	consumeAndSendToAggregator, err := NewKafkaConsume(kafkaTopic, r, c)
	if err != nil {
		log.Fatal(err)
	}
	consumeAndSendToAggregator.start()
}
