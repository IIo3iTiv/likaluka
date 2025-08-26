package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	addr := os.Getenv("AGREGATOR_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}
	log.Println("Addr:", addr)
	// Инициализация stream sender
	sender, err := NewStreamSender(addr)
	if err != nil {
		log.Fatalf("Failed to create sender: %v", err)
	}
	defer func() {
		receivedCount, err := sender.Close()
		if err != nil {
			log.Fatalf("Stream closed with error: %v", err)
		}
		log.Printf("Server confirmed receipt of %d rows", receivedCount)
	}()

	// Чтение и отправка CSV
	log.Printf("Starting CSV processing...")

	if err := ReadAndStreamCSVFast(filepath.Join(GetExecuteFilePath(), "datasets", "current_1.csv"), 1, sender); err != nil {
		log.Fatalf("Failed to process CSV: %v", err)
	}
}

func GetExecuteFilePath() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}
