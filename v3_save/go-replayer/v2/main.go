package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	pb "itc/proto/v2"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StreamSender struct {
	client pb.DataServiceClient
	stream pb.DataService_StreamDataClient
	ctx    context.Context
	cancel context.CancelFunc
}

func NewStreamSender(addr string) (*StreamSender, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(30*1024*1024),
			grpc.MaxCallSendMsgSize(30*1024*1024),
		))
	if err != nil {
		return nil, err
	}

	client := pb.NewDataServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())

	stream, err := client.StreamData(ctx)
	if err != nil {
		cancel()
		conn.Close()
		return nil, err
	}

	return &StreamSender{
		client: client,
		stream: stream,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (ss *StreamSender) SendRow(guid uint32, timestamp int64, r, s, t float32) error {
	row := &pb.Request{
		Guid:      guid,
		Timestamp: timestamp,
		R:         r,
		S:         s,
		T:         t,
	}
	return ss.stream.Send(row)
}

func (ss *StreamSender) Close() (uint64, error) {
	response, err := ss.stream.CloseAndRecv()
	if err != nil {
		return 0, err
	}

	if !response.Success {
		return response.ReceivedCount, fmt.Errorf("server error: %s", response.Msg)
	}

	return response.ReceivedCount, nil
}

func readAndStreamCSV(filename string, guid uuid.UUID, sender *StreamSender) error {
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
func readAndStreamCSVFast(filename string, guid uint32, sender *StreamSender) error {
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

	fr, fs, ft, err := findColumnRST(file)
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
			Guid:      guid,
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

func findColumnRST(file *os.File) (fr, fs, ft int, err error) {
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

	if err := readAndStreamCSV(filepath.Join(GetExecuteFilePath(), "datasets", "current_1.csv"), uuid.New(), sender); err != nil {
		log.Fatalf("Failed to process CSV: %v", err)
	}
}

func GetExecuteFilePath() string {
	ex, _ := os.Executable()
	return filepath.Dir(ex)
}
