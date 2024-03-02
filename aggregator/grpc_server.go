package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/adalbertjnr/ws-person/types"
	"google.golang.org/grpc"
)

type GRPCServerTransport struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCServerTransport(s Aggregator) *GRPCServerTransport {
	return &GRPCServerTransport{
		svc: s,
	}
}

func GRPCServer(grpcServerAddr string, c Aggregator) error {
	lis, err := net.Listen("tcp", grpcServerAddr)
	if err != nil {
		return fmt.Errorf("failed to listen %w", err)
	}
	s := grpc.NewServer()
	types.RegisterAggregatorServer(s, NewGRPCServerTransport(c))
	log.Println("grpc server listening at", grpcServerAddr)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve the grpc server %w", err)
	}
	return nil
}

func (g *GRPCServerTransport) Aggregate(ctx context.Context, p *types.AggregatePerson) (*types.None, error) {
	const ReceivingInGRPCServer = "receiving in grpc server"
	d := types.Person{
		Id:    int(p.Id),
		Name:  p.Name,
		Age:   uint8(p.Age),
		Role:  p.Role,
		Stage: ReceivingInGRPCServer,
	}
	return &types.None{}, g.svc.Insert(d)
}
