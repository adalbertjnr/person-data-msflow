package client

import (
	"context"

	"github.com/adalbertjnr/ws-person/types"
)

type ClientPicker interface {
	Aggregate(context.Context, *types.AggregatePerson) error
	GetPersonById(context.Context, int) (*types.Person, error)
}
