package client

import (
	"context"
	"fmt"

	"github.com/adalbertjnr/ws-person/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClientEndpoint struct {
	endpoint string
	types.AggregatorClient
}

func NewGRPCClientEndpoint(e string) (*GRPCClientEndpoint, error) {
	conn, err := grpc.Dial(e, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("error dialing grpc server %w", err)
	}
	c := types.NewAggregatorClient(conn)
	return &GRPCClientEndpoint{
		endpoint:         e,
		AggregatorClient: c,
	}, nil
}

func (g *GRPCClientEndpoint) Aggregate(ctx context.Context, data *types.AggregatePerson) error {
	_, err := g.AggregatorClient.Aggregate(ctx, data)
	return err
}
