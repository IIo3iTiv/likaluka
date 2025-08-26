package main

import (
	"log"
	"os"
	"path/filepath"
)

var g GrpcBatchSender

func main() {
	var err error
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	addr := os.Getenv("AGREGATOR_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}
	log.Println("Addr:", addr)

	g, err = InitGRPC(addr)
	if err != nil {
		log.Printf("Error: InitGRPC: %s", err.Error())
		return
	}
	log.Println("Succes. InitGRPC")

	InitServer()
}

func GetExecuteFilePath() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}
