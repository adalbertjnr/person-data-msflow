package main

import (
	"github.com/adalbertjnr/ws-person/types"
	"github.com/sirupsen/logrus"
)

type DataMiddlewareLogger struct {
	next Replacer
	log  *logrus.Logger
}

func NewDataMiddlewareLogger(r Replacer) *DataMiddlewareLogger {
	return &DataMiddlewareLogger{
		next: r,
		log:  types.LoggerJSONClient(),
	}
}

func (d *DataMiddlewareLogger) ReplaceData(data types.Person) *types.Person {
	defer func() {
		d.log.WithFields(logrus.Fields{
			"id":    data.Id,
			"name":  data.Name,
			"age":   data.Age,
			"role":  data.Role,
			"stage": data.Stage,
		}).Infoln("consuming from kafka and sending to aggregator")
	}()
	return d.next.ReplaceData(data)
}
