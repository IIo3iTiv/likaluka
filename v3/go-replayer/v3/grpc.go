package main

import (
	"context"
	pb "itc/proto/v4"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcBatchSender struct {
	client pb.DataServiceClient
}

func InitGRPC(addr string) (g GrpcBatchSender, err error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10*1024*1024),
			grpc.MaxCallSendMsgSize(5*1024*1024),
		),
	)
	if err != nil {
		return
	}
	g.client = pb.NewDataServiceClient(conn)
	return
}

func (g *GrpcBatchSender) SendDataBatch(batch *pb.DataBatch) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := g.client.SendDataBatch(ctx, batch)
	if err != nil {
		log.Printf("ID: %d. Error: %s", batch.GetGuid(), err.Error())
	}
	log.Printf("ID: %d. Success. Message: %s", batch.Guid, response.GetMsg())
}
