package client

import (
	"context"

	"github.com/adalbertjnr/ws-person/types"
)

type ClientPicker interface {
	Aggregate(context.Context, *types.AggregatePerson) error
}
