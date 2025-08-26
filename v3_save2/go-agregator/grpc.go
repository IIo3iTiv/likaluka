package main

import (
	"context"
	"itc/proto"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var datasets map[string]*proto.DataMessage
var count uint = 0
var t1 time.Time
var t2 time.Time
var count2 uint = 0

type server struct {
	proto.UnimplementedDataServiceServer
}

func InitGRPC() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err.Error())
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor))
	proto.RegisterDataServiceServer(s, &server{})
	log.Print("Server running")
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to server: %s", err.Error())
	}
}

func (s *server) Send(ctx context.Context, req *proto.DataMessage) (*emptypb.Empty, error) {
	// datasets[req.Guid] = req
	// log.Println("ALALLAALAL", req.String())
	// k.Publish(fmt.Sprintf("R=%f;S=%f;T=%f;GUID=%s;Date=%d;", req.GetR(), req.GetS(), req.GetT(), req.GetGuid(), req.GetDate()))
	if count == 0 {
		t1 = time.Unix(req.Date, 0)
		t2 = time.Now()
	}
	if count2 >= 25600 {
		log.Println("Всего файлов: ", count)
		log.Println("T2 ", t2.Format(time.RFC3339Nano))
		log.Println("Разница: ", time.Since(t2).Seconds())
		count2 = 0
		t2 = time.Now()
	}
	if req.Guid == "lalka" {
		log.Println("The End")
		log.Println("Всего файлов: ", count)
		log.Println("T1 ", t1.Format(time.RFC3339Nano))
		log.Println("TimeNow: ", time.Now().Format(time.RFC3339Nano))
		log.Println("Разница: ", time.Since(t1).Minutes())
	}
	count++
	count2++
	return &emptypb.Empty{}, nil
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// log.Printf("Received request: %v", info.FullMethod)
	return handler(ctx, req)
}
