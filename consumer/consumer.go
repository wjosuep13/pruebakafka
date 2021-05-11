package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

const (
	topic         = "message-log"
	brokerAddress = "localhost:19092"
)

type Message struct {
	Name         string `json:"name"`
	Location     string `json:"location"`
	Age          int    `json:"age"`
	Infectedtype string `json:"infectedtype"`
	State        string `json:"state"`
	Ruta         string `json:"ruta"`
}

func listener() {

	ctx := context.Background()

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

		resp, err := http.Post(host, "application/json; charset=utf-8", bytes.NewBuffer(msg.Value))
		if err != nil {
			log.Fatalln(err)
			return
		}

		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		// Convert response body to string
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}

	return
}

func main() {
	listener()
}
