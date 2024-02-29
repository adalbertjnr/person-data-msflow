package main

import (
	"fmt"

	"github.com/adalbertjnr/ws-person/types"
)

type Aggregator interface {
	// Get person from the aggregator database
	Get(id int) (types.Person, error)
	// Insert person into aggregator database
	Insert(data types.Person) error
}

type DataStore struct {
	person map[int]types.Person
}

func NewDataStore() *DataStore {
	return &DataStore{
		person: make(map[int]types.Person),
	}
}

func (d *DataStore) Insert(data types.Person) error {
	fmt.Printf("storing %v\n", data)
	d.person[data.Id] = data
	return nil
}

func (d *DataStore) Get(id int) (types.Person, error) {
	return d.person[id], nil
}
