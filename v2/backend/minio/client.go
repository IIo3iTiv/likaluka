package minio

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type client struct {
	mc *minio.Client
}

var c client

func Init(ctx context.Context, endpoint, u, p, bucket string, ssl bool) (err error) {
	c.mc, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(u, p, ""),
		Secure: ssl,
	})
	if err != nil {
		return
	}
	if !Check() {
		err = fmt.Errorf("Minio: client not active")
		return
	}

	return
}

func Check() bool {
	return c.mc.IsOnline()
}
