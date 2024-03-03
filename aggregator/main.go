package main

import (
	"log"
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
	log.Fatal(HTTPServer(httpServerAddr))
}
