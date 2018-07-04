package main

import (
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	port := os.Getenv("port")
	if port == "" {
		port = "1337"
	}
	http.HandleFunc("/feed", liveFeed)
	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func liveFeed(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS: Could not establish websocket connection. Error: ", err)
		return
	}
	log.Println("WS: Connection established.")
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "178.128.173.16"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("WS: There was a problem reading the message. Error: ", err)
			return
		}
		topic := "nyc_driver_loc_data"
		fMsg := string(message)
		fmt.Println(fMsg)

		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(nil),
			Value:          message,
		}, nil)
	}
}
