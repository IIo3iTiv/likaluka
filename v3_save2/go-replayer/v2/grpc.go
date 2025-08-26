package main

import (
	"context"
	"fmt"
	pb "itc/proto/v2"
	"log"
	"time"

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
	var maxRetry = 30
	var timeSleep = time.Second * 20
	var count = 0

	var fSleep = func(ts time.Duration) error {
		count++
		if maxRetry == count {
			return fmt.Errorf("%s", "the maximum number of connection attempts has been reached")
		}
		time.Sleep(timeSleep)
		return nil
	}

	for {
		conn, err := grpc.NewClient(addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(30*1024*1024),
				grpc.MaxCallSendMsgSize(30*1024*1024),
			))
		if err != nil {
			log.Printf("GRPC: NewClient: %s", err.Error())
			if err := fSleep(timeSleep); err != nil {
				// conn.Close()
				return nil, err
			}
		}

		client := pb.NewDataServiceClient(conn)
		ctx, cancel := context.WithCancel(context.Background())

		stream, err := client.StreamData(ctx)
		if err != nil {
			cancel()
			conn.Close()
			log.Printf("StreamData: %s", err.Error())
			if err := fSleep(timeSleep); err != nil {
				return nil, err
			}
		}
		if err == nil {
			return &StreamSender{
				client: client,
				stream: stream,
				ctx:    ctx,
				cancel: cancel,
			}, nil
		}
	}
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
