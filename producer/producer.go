package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/oklog/ulid/v2"

	"github.com/gorilla/mux"

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

func test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, "ok")

}

func publish(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	reqBody, _ := ioutil.ReadAll(r.Body)
	msg := Message{Ruta: "Kafka"}
	json.Unmarshal(reqBody, &msg)
	data, err := json.Marshal(msg)
	ctx := context.Background()
	if err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	produce(ctx, data)
	// Obtener el mensaje enviado desde la forma

}

func produce(ctx context.Context, data []byte) {
	// initialize a counter

	l := log.New(os.Stdout, "kafka writer: ", 0)
	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		// assign the logger to the writer
		Logger: l,
	})

	// each kafka message has a key and value. The key is used
	// to decide which partition (and consequently, which broker)
	// the message gets published on
	err := w.WriteMessages(ctx, kafka.Message{
		Key: []byte(Ulid()),
		// create an arbitrary message payload for the value
		Value: []byte(string(data)),
	})
	if err != nil {
		panic("could not write message " + err.Error())
	}

}

func Ulid() string {
	t := time.Now().UTC()
	id := ulid.MustNew(ulid.Timestamp(t), rand.Reader)

	return id.String()
}

func main() {
	handleRequests()
}

func handleRequests() {
	port := "5500"
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", test).Methods(http.MethodGet)
	r.HandleFunc("/data", publish).Methods(http.MethodPost)
	r.Use(mux.CORSMethodMiddleware(r))

	http.ListenAndServe(":"+port, r)
}
