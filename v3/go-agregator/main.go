package main

import (
	"log"
	"os"
	"path/filepath"
)

var k Kafka

func main() {
	var err error
	k, err = InitKafka()
	if err != nil {
		log.Printf("Error: InitKafka: %s", err.Error())
		return
	}
	log.Println("1.", 100%1000)
	log.Println("2.", 10000%1000)
	log.Println("3.", 10000%1000)

	log.Println("Success InitKafka")
	InitGRPC()
}

func GetExecuteFilePath() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}
