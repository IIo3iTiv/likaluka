package main

import (
	"context"
	"itc/proto/v1"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	conn   *grpc.ClientConn
	client proto.DataServiceClient
}

func InitGRPCClient(addr string) (g GrpcClient, err error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}

	g.conn = conn
	g.client = proto.NewDataServiceClient(g.conn)
	return
}

func (grpc *GrpcClient) Send(ctx context.Context, r Record, guid string, date time.Time) (err error) {
	_, err = grpc.client.Send(ctx, &proto.DataMessage{
		R:    r.R,
		S:    r.S,
		T:    r.T,
		Guid: guid,
		Date: date.UnixNano(),
	})
	return err
}
