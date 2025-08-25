package main

import (
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"

	pb "itc/proto/v2"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDataServiceServer
	totalReceived uint64
}

func (s *server) StreamData(stream pb.DataService_StreamDataServer) error {
	var receivedCount uint64 = 0

	for {
		row, err := stream.Recv()
		if err == io.EOF {
			// Клиент завершил отправку
			total := atomic.AddUint64(&s.totalReceived, receivedCount)
			log.Printf("Stream completed. Received %d rows in this stream, total: %d", receivedCount, total)
			return stream.SendAndClose(&pb.Response{
				ReceivedCount: receivedCount,
				Success:       true,
				Msg:           "All data received successfully",
			})
		}
		if err != nil {
			return err
		}

		// Обработка данных с учетом новых полей
		s.processRow(row)

		receivedCount++

		// Периодическое логирование
		if receivedCount%10000 == 0 {
			log.Printf("Received %d rows (GUID: %d, Timestamp: %d)", receivedCount, row.Guid, row.Timestamp)
		}
	}
}

var count uint = 0
var t1 time.Time
var t2 time.Time
var count2 uint = 0

func (s *server) processRow(row *pb.Request) {
	// Ваша логика обработки данных
	// Теперь есть доступ к row.Guid и row.Timestamp
	// Пример:
	//timestamp := time.Unix(row.Timestamp, 0)
	//log.Printf("Processing: GUID=%s, Time=%s, R=%.32f, S=%.32f, T=%.32cf", row.Guid, timestamp.Format(time.RFC3339), row.R, row.S, row.T)
	if count == 0 {
		t1 = time.Unix(row.Timestamp, 0)
		t2 = time.Now()
	}
	if count2 >= 25600 {
		log.Println("Всего файлов: ", count)
		log.Println("T2 ", t2.Format(time.RFC3339Nano))
		log.Println("Разница: ", time.Since(t2).Seconds())
		count2 = 0
		t2 = time.Now()
	}
	count++
	count2++

	// Можно отправлять в канал для параллельной обработки
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(20*1024*1024),
		grpc.MaxSendMsgSize(20*1024*1024),
	)

	service := &server{}
	pb.RegisterDataServiceServer(s, service)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	_ = t1
}
