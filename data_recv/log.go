package main

import (
	"github.com/adalbertjnr/ws-person/types"
	"github.com/sirupsen/logrus"
)

type DataMiddlewareLogger struct {
	next Producer
	log  *logrus.Logger
}

func NewDataMiddlewareLogger(n Producer) *DataMiddlewareLogger {
	return &DataMiddlewareLogger{
		next: n,
		log:  types.LoggerJSONClient(),
	}
}

func (d *DataMiddlewareLogger) ProduceToKafka(data types.Person) error {
	defer func() {
		d.log.WithFields(logrus.Fields{
			"id":    data.Id,
			"name":  data.Name,
			"age":   data.Age,
			"role":  data.Role,
			"stage": data.Stage,
		}).Infoln("producing to kafka")
	}()
	return d.next.ProduceToKafka(data)
}
