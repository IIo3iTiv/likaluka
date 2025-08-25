package main

import (
	"context"
	"fmt"
	pb "itc/proto/v3"
	"log"
	"net"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDataServiceServer
}

func InitGrpc() (err error) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		err = fmt.Errorf("listen: %s", err.Error())
		return
	}
	s := grpc.NewServer()
	pb.RegisterDataServiceServer(s, &server{})

	err = s.Serve(lis)
	if err != nil {
		err = fmt.Errorf("server: %s", err.Error())
	}
	return err
}

func (s *server) SendDataBatch(ctx context.Context, batch *pb.DataBatch) (resp *pb.Response, err error) {
	log.Printf("Received batch. Guid: %d. Length: %d", batch.GetGuid(), len(batch.Points))
	var dataJSON DataJSON
	dataJSON.Guid = batch.GetGuid()
	for _, bp := range batch.Points {
		dataJSON.Data = append(dataJSON.Data, BatchJSON{
			Date: bp.Date.AsTime().UnixNano(),
			R:    bp.R,
			S:    bp.S,
			T:    bp.T,
		})
	}

	json := jsoniter.ConfigFastest
	jsonData, err := json.MarshalToString(&dataJSON)
	if err != nil {
		resp = &pb.Response{Success: false, Msg: fmt.Sprintf("Error: JSONMarshal: %s", err.Error())}
		log.Printf("Error. MershalJSON. Guid: %d. ErrorMsg: %s", batch.GetGuid(), err.Error())
		return
	}
	err = KafkaPublish(jsonData)
	if err != nil {
		resp = &pb.Response{Success: false, Msg: fmt.Sprintf("Error: KafkaPublish: %s", err.Error())}
		log.Printf("Error. KafkaPublish. ErrorMsg: %s", err.Error())
		return
	}
	resp = &pb.Response{Success: true, Msg: "Ok"}
	return
}
