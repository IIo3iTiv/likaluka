package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/jszwec/csvutil"
)

type Record struct {
	R float64 `csv:"current_R"`
	S float64 `csv:"current_S"`
	T float64 `csv:"current_T"`
}

func ReadCSV(grpc *GrpcClient, file *os.File, data chan Record) (rs []Record, err error) {
	dec, err := csvutil.NewDecoder(csv.NewReader(file))
	if err != nil {
		err = fmt.Errorf("csvutil.NewDecoder:%s", err.Error())
		return
	}
	for {
		var r Record
		err = dec.Decode(&r)
		if err != nil {
			log.Println("ReadCSV The End")
			// r = Record{R: 9, S: 9, T: 9}
			// grpc.Send(context.Background(), r, uuid.New().String(), time.Now())
			// close(data)
			return
		}
		rs = append(rs, r)
		// data <- r
		// grpc.Send(context.Background(), r, uuid.New().String(), time.Now())
	}
}
