package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func main() {
	addr := os.Getenv("AGREGATOR_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}
	log.Println("Addr:", addr)

	var grpcc GrpcClient
	var err error
	maxAttempt := 30
	sleep := time.Second * 5
	for {
		maxAttempt--
		if maxAttempt == 0 {
			log.Printf("ERROR: GRPC failed to connect")
			return
		}
		grpcc, err = InitGRPCClient(addr)
		if err != nil {
			log.Printf("ERROR: InitGRPCClient:%s", err.Error())
			time.Sleep(sleep)
			continue
		}
		break
	}
	err = ReadDataset(&grpcc)
	if err != nil {
		log.Printf("Error: ReadDataset: %s", err.Error())
	}
	for {
	}
}

func ReadDataset(grpc *GrpcClient) (err error) {
	datasetPath := filepath.Join(GetExecuteFilePath(), "datasets", "current_1.csv")
	dataset1, err := os.Open(datasetPath)
	if err != nil {
		return
	}
	log.Printf("ReadCSV:%s", datasetPath)

	// ctx := context.Background()
	data := make(chan Record)
	rs, err := ReadCSV(grpc, dataset1, data)
	if err != nil && err != io.EOF {
		log.Println("ReadCSV: ", err.Error())
	}
	ctx := context.Background()
	guid := uuid.New().String()
	for _, r := range rs {
		grpc.Send(ctx, r, guid, time.Now())
	}
	grpc.Send(ctx, Record{R: 0, S: 0, T: 0}, "lalka", time.Now())
	// for d := range data {
	// 	if d.R == 9 {
	// 		if d.S == 9 && d.T == 9 {
	// 			grpc.Send(ctx, d, "lalka", time.Now())
	// 		}
	// 	}
	// 	grpc.Send(ctx, d, uuid.New().String(), time.Now())
	// }
	log.Println("Success read CSV")
	return
}

func GetExecuteFilePath() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}
