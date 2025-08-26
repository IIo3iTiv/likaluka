package main

import (
	"log"
	"time"
)

func InitService(name string, x func() error) {
	var maxRetry uint16 = 30
	var timeSleep time.Duration = 10 * time.Second
	var count uint16 = 0

	log.Printf("Init Start. Service: %s", name)
	for {
		log.Printf("Service: %s. Retry: %d", name, count)
		err := x()
		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
		if err == nil {
			break
		}
		if maxRetry == count {
			log.Fatalf("Error: maxRetry")
			return
		}
		time.Sleep(timeSleep)
		count++
	}
	log.Printf("Success. Service: %s", name)
}

func main() {
	InitService("Kafka", InitKafka)
	InitService("GRPC", InitGrpc)
}
