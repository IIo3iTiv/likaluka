package main

import (
	"log"
	"time"
)

func InitService(x func() error) {
	var maxRetry uint16 = 30
	var timeSleep time.Duration = 10 * time.Second
	var count uint16 = 0

	log.Println("Init Start")
	for {
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
	log.Printf("Success")
}

func main() {
	InitService(InitKafka)
	InitService(InitGrpc)
}
