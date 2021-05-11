package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/segmentio/kafka-go"
)

const (
	topic         = "messages"
	brokerAddress = "localhost:9092"
)

type Message struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	VaccineType string `json:"vaccine_type"`
	Ruta        string `json:"ruta"`
}

type Redis struct {
	DB   string  `json:"db"`
	Data Message `json:"data"`
}

func listener() {

	ctx := context.Background()
	host := "http://35.208.191.6:8080/add"
	l := log.New(os.Stdout, "kafka reader: ", 0)
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		GroupID: "my-group",
		// assign the logger to the reader
		Logger: l,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
		var jeison Message
		json.Unmarshal(msg.Value, &jeison)

		redis := Redis{DB: "covid", Data: jeison}

		rbtes, err := json.Marshal(redis)

		resp, err := http.Post(host, "application/json; charset=utf-8", bytes.NewBuffer(rbtes))
		if err != nil {
			log.Fatalln(err)
			fmt.Println("error perrillo")
			return
		}
		defer resp.Body.Close()

	}

}

func main() {
	listener()
}
