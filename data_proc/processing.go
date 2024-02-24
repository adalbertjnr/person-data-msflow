package main

import "github.com/adalbertjnr/ws-person/types"

type Replacer interface {
	ReplaceData(data types.Person) (newData types.Person)
}

type Replace struct {
	currentStage string
}

func NewReplace() *Replace {
	return &Replace{}
}

func (r *Replace) ReplaceData(data types.Person) (newData types.Person) {
	r.currentStage = "data_proc stage"
	newData = data
	newData.Stage = r.currentStage
	return newData
}
