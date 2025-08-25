package main

import (
	"bufio"
	"fmt"
	pb "itc/proto/v3"
	"os"
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CSVRecord struct {
	R float32
	S float32
	T float32
}

// uint16 = 2byte - guid = 2 byte
// float32 = 4byte - r, s ,t = 12 byte
// int64 = 8byte - timestamp = 8 byte
// = 22 byte
func (g *GrpcBatchSender) ReadCSV(filename string, guid uint32) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fr, fs, ft, err := FindColumnRST(file)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	batch := make([]CSVRecord, 0, 256)

	if scanner.Scan() {

	}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		r, _ := strconv.ParseFloat(parts[fr], 32)
		s, _ := strconv.ParseFloat(parts[fs], 32)
		t, _ := strconv.ParseFloat(parts[ft], 32)
		row := CSVRecord{
			R: float32(r),
			S: float32(s),
			T: float32(t),
		}
		batch = append(batch, row)
		if len(batch) == 256 {
			go g.SendData(guid, batch)
			batch = batch[:0]
		}
	}
	go g.SendData(guid, batch)
	return nil
}

func (g *GrpcBatchSender) SendData(guid uint32, data []CSVRecord) {
	batch := pb.DataBatch{}
	batch.Guid = guid
	for _, d := range data {
		batch.Points = append(batch.Points, &pb.DataPoint{
			Date: timestamppb.Now(),
			R:    d.R,
			S:    d.S,
			T:    d.T,
		})
	}
	g.SendDataBatch(&batch)
}

func FindColumnRST(file *os.File) (fr, fs, ft int, err error) {
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "R") && strings.Contains(line, "S") && strings.Contains(line, "T") {
			lines := strings.Split(line, ",")
			for i, l := range lines {
				if strings.Contains(l, "R") {
					fr = i
					continue
				}
				if strings.Contains(l, "S") {
					fs = i
					continue
				}
				if strings.Contains(l, "T") {
					ft = i
					continue
				}
			}
		}
	}
	if fs+fr+ft != 3 {
		err = fmt.Errorf("%s", "not found fields value")
	}
	return
}
