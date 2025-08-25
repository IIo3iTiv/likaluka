package main

import (
	"log"
)

var k Kafka

func main() {
	var err error
	k, err = InitKafka()
	if err != nil {
		log.Printf("Error: Kafka: %s", err.Error())
	}
	log.Println("Success init kafka")

	_, err = InitGRPC(":50051")
	if err != nil {
		log.Printf("Error: InitGRPC: %s", err.Error())
	}
}
