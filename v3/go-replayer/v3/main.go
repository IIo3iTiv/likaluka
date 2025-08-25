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

	g, err := InitGRPC(addr)
	if err != nil {
		log.Printf("Error: InitGRPC: %s", err.Error())
		return
	}
	log.Println("Succes. InitGRPC")

	err = g.ReadCSV(filepath.Join(GetExecuteFilePath(), "datasets", "current_1.csv"), uint32(1))
	if err != nil {
		log.Printf("Error: ReadAndStremCSVFast: %s", err.Error())
	}
	log.Println("Success. ReadCSVFast")
}

func GetExecuteFilePath() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}
