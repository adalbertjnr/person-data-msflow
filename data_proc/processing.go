package main

import "github.com/adalbertjnr/ws-person/types"

type Replacer interface {
	ReplaceData(data types.Person) *types.Person
}

type Replace struct {
	currentStage string
}

func NewReplace() *Replace {
	return &Replace{}
}

func (r *Replace) ReplaceData(data types.Person) *types.Person {
	r.currentStage = currentStage
	data.Stage = r.currentStage
	return &data
}
