package main

import (
	"context"
	"log"

	"github.com/adalbertjnr/ws-person/aggregator/client"
	"github.com/adalbertjnr/ws-person/types"
)

const (
	httpServerAddr = ":3001"
	grpcServerAddr = ":4001"
)

func main() {
	go func() {
		if err := GRPCServer(grpcServerAddr, NewDataStore()); err != nil {
			log.Fatal(err)
		}
	}()
	c, err := client.NewGRPCClientEndpoint(grpcServerAddr)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Aggregate(context.Background(), &types.AggregatePerson{
		Id:    int32(2),
		Name:  "foo",
		Age:   int32(20),
		Role:  "barRole",
		Stage: "checking grpc client",
	}); err != nil {
		log.Fatal(err)
	}
	log.Fatal(HTTPServer(httpServerAddr))
}
