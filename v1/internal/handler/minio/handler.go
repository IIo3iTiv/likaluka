package minio

import "miniogo/pkg/minio"

type Handler struct {
	minioService minio.Client
}

func NewMinioHandler(
	minioService minio.Client,
) *Handler {
	return &Handler{
		minioService: minioService,
	}
}
