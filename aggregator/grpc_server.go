package main

import (
	"fmt"
	"log"
	"net"

	"github.com/adalbertjnr/ws-person/types"
	"google.golang.org/grpc"
)

type GRPCServerTransport struct {
	types.UnimplementedAggregatorServer
}

func GRPCServer(grpcServerAddr string) error {
	lis, err := net.Listen("tcp", grpcServerAddr)
	if err != nil {
		return fmt.Errorf("failed to listen %w", err)
	}
	s := grpc.NewServer()
	types.RegisterAggregatorServer(s, GRPCServerTransport{})
	log.Println("grpc server listening at", grpcServerAddr)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve the grpc server %w", err)
	}
	return nil
}
