package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	pb "itc/proto/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ReadAndStreamCSV(filename string, guid uuid.UUID, sender *StreamSender) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Пропускаем заголовок если есть
	if _, err := reader.Read(); err != nil && err != io.EOF {
		return err
	}

	sentCount := 0
	startTime := time.Now()
	lastLogTime := startTime

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Парсим поля (5 полей: guid, timestamp, r, s, t)
		if len(record) != 3 {
			return fmt.Errorf("invalid record length: %d", len(record))
		}

		// guid := record[0]

		// timestamp, err := strconv.ParseInt(record[1], 10, 64)
		// if err != nil {
		// 	return fmt.Errorf("invalid timestamp: %v", err)
		// }

		r, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return err
		}
		_ = r
		s, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return err
		}
		_ = s
		t, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return err
		}
		_ = t

		// if err := sender.SendRow(guid.String(), time.Now().Unix(), r, s, t); err != nil {
		// 	return err
		// }

		sentCount++

		// Контроль скорости: 25600/сек = 256 каждые 10ms
		if sentCount%256 == 0 {
			time.Sleep(time.Millisecond * 10)

			// Логирование скорости каждые 0.1 секунды
			currentTime := time.Now()
			if currentTime.Sub(lastLogTime) >= 100*time.Millisecond {
				elapsed := currentTime.Sub(startTime).Seconds()
				rate := float64(sentCount) / elapsed
				log.Printf("Sent %d rows, current rate: %.0f rows/sec", sentCount, rate)
				lastLogTime = currentTime
			}
		}
	}

	// Финальное логирование
	elapsed := time.Since(startTime).Seconds()
	rate := float64(sentCount) / elapsed
	log.Printf("Finished: %d rows in %.2f seconds (%.0f rows/sec)",
		sentCount, elapsed, rate)

	return nil
}

// uint16 = 2byte - guid = 2 byte
// float32 = 4byte - r, s ,t = 12 byte
// int64 = 8byte - timestamp = 8 byte
// = 22 byte
func ReadAndStreamCSVFast(filename string, guid uint32, sender *StreamSender) error {
	send := func(ds []*pb.Request, sender *StreamSender) {
		for _, d := range ds {
			sender.SendRow(d.Guid, d.Timestamp, d.R, d.S, d.T)
		}
	}

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
	buf := make([]byte, 0, 256*1024)
	scanner.Buffer(buf, 10*1024*1024)
	sendCount := 0
	batch := make([]*pb.Request, 0, 256)

	if scanner.Scan() {

	}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		r, _ := strconv.ParseFloat(parts[fr], 32)
		s, _ := strconv.ParseFloat(parts[fs], 32)
		t, _ := strconv.ParseFloat(parts[ft], 32)
		row := &pb.Request{
			Guid:      uint32(sendCount),
			Timestamp: time.Now().Unix(),
			R:         float32(r),
			S:         float32(s),
			T:         float32(t),
		}
		batch = append(batch, row)
		sendCount++

		if len(batch) >= 256 {
			go send(batch, sender)
			batch = batch[:0]
		}
	}
	go send(batch, sender)
	return nil
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
