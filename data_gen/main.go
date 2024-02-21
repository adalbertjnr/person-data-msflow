package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/adalbertjnr/ws-person/types"
	"github.com/gorilla/websocket"
)

const (
	SendingFromGeneratorToReceiver = "mockDataProducer"
	BatchSize                      = 10
	endpoint                       = "ws://127.0.0.1:3000"
)

func generatePersonData(batchSize int) []types.Person {
	newBatchPersonData := []types.Person{}
	for i := 0; i < batchSize; i++ {
		name, role := pickRandomNameAndJob(personMockData())
		person := types.Person{
			Name:  name,
			Age:   uint8(rand.Intn(100)),
			Role:  role,
			Stage: SendingFromGeneratorToReceiver,
		}
		newBatchPersonData = append(newBatchPersonData, person)
	}
	return newBatchPersonData
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		log.Fatalf("error dialing with ws client to %s: %v", endpoint, err)
	}
	ticker := time.NewTicker(time.Second * 10)

	for {
		if err := conn.WriteJSON(generatePersonData(BatchSize)); err != nil {
			fmt.Printf("error writing json from the ws client %v", err)
		}
		<-ticker.C
	}
}

func pickRandomNameAndJob(data map[string]string) (string, string) {
	var saveName, saveJob []string
	for dataName, dataJob := range data {
		saveName = append(saveName, dataName)
		saveJob = append(saveJob, dataJob)
	}
	randomPicker := rand.Intn(len(data))
	return saveName[randomPicker], saveJob[randomPicker]
}
