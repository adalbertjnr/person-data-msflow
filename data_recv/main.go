package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adalbertjnr/ws-person/types"
	"github.com/gorilla/websocket"
)

type DataReceiver struct {
	conn *websocket.Conn
	kp   Producer
}

const (
	wsServerPort = ":3000"
	kafkaTopic   = "wstopic"
	stageStatus  = "data_receiver sending to kafka"
)

func main() {

	rec, err := NewDataReceiver(kafkaTopic)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", rec.wsHandler)

	fmt.Println("websocket server started on port", wsServerPort)
	log.Fatal(http.ListenAndServe(wsServerPort, nil))
}

func NewDataReceiver(topic string) (*DataReceiver, error) {
	var (
		kafkaProducerClient Producer
		err                 error
	)
	kafkaProducerClient, err = NewKafkaProduce(topic)
	if err != nil {
		return nil, err
	}

	kafkaProducerClient = NewDataMiddlewareLogger(kafkaProducerClient)
	return &DataReceiver{
		kp: kafkaProducerClient,
	}, nil
}

func (d *DataReceiver) wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	d.conn = conn

	go d.startReceiveData()
}

func (d *DataReceiver) startReceiveData() {
	fmt.Println("new websocket client connected")
	for {
		var dataSlice []types.Person
		if err := d.conn.ReadJSON(&dataSlice); err != nil {
			log.Println("read json in ws server error:", err)
		}
		for _, data := range dataSlice {
			data.Stage = stageStatus
			if err := d.kp.ProduceToKafka(data); err != nil {
				log.Println("producing to kafka in ws server error:", err)
			}
		}
	}
}
