package main

import (
	"log"
)

const kafkaTopic = "wstopic"

func main() {
	c, err := NewKafkaConsume(kafkaTopic, NewReplace())
	if err != nil {
		log.Fatal(err)
	}
	c.start()
}
