package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/adalbertjnr/ws-person/types"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	SendingFromGeneratorToReceiver = "mockDataProducer"
	BatchSize                      = 10
	endpoint                       = "ws://127.0.0.1:3000"
)

func generatePersonData(batchSize int, l *logrus.Logger) []types.Person {
	newBatchPersonData := []types.Person{}
	for i := 0; i < batchSize; i++ {
		name, role := pickRandomNameAndJob(personMockData())
		person := types.Person{
			Id:    int(idGenerator()),
			Name:  name,
			Age:   uint8(rand.Intn(100)),
			Role:  role,
			Stage: SendingFromGeneratorToReceiver,
		}
		l.WithFields(logrus.Fields{"id": person.Id, "name": person.Name, "age": person.Age, "role": person.Role, "stage": person.Stage}).
			Info("generating data")
		newBatchPersonData = append(newBatchPersonData, person)
	}
	return newBatchPersonData
}

func idGenerator() int64 {
	return int64(rand.Intn(100000))
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
func main() {
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		log.Fatalf("error dialing with ws client to %s: %v", endpoint, err)
	}
	ticker := time.NewTicker(time.Second * 10)

	for {
		if err := conn.WriteJSON(generatePersonData(BatchSize, types.LoggerJSONClient())); err != nil {
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
