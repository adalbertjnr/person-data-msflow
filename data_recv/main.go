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
	kp   *KafkaProduce
}

const (
	wsServerPort = ":3000"
	kafkaTopic   = "mytopic"
)

func main() {
	rec, err := NewDataReceiver(kafkaTopic)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", httpHandlerWrapper(rec.wsHandler))

	fmt.Println("websocket server started on port", wsServerPort)
	log.Fatal(http.ListenAndServe(wsServerPort, nil))
}

func NewDataReceiver(topic string) (*DataReceiver, error) {
	kafkaProducerClient, err := NewKafkaProduce(topic)
	if err != nil {
		return nil, err
	}

	return &DataReceiver{
		kp: kafkaProducerClient,
	}, nil
}

func (d *DataReceiver) wsHandler(w http.ResponseWriter, r *http.Request) error {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	d.conn = conn

	errCh := make(chan error)
	go func() {
		errCh <- d.startReceiveData()
	}()

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}

func (d *DataReceiver) startReceiveData() error {
	fmt.Println("new websocket client connected")
	for {
		var dataSlice []types.Person
		if err := d.conn.ReadJSON(&dataSlice); err != nil {
			return err
		}
		for _, data := range dataSlice {
			if err := d.kp.ProduceToKafka(data); err != nil {
				return err
			}
		}
	}
}

type fnType func(w http.ResponseWriter, r *http.Request) error

func httpHandlerWrapper(fn fnType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			log.Printf("error from the websocket server handler %v", err)
		}
	}
}
