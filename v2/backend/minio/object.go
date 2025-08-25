package minio

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type Object struct {
	LocalPath  string `json:"FilePath"`
	LocalName  string `json:"LocalName"`
	Name       string `json:"FileName"`
	UploadInfo Object_UploadInfo
}

type Object_UploadInfo struct {
	Bucket            string    `json:"Bucket"`
	Key               string    `json:"Key"`
	ETag              string    `json:"ETag"`
	Size              int64     `json:"Size"`
	LastModified      time.Time `json:"LastModified"`
	Location          string    `json:"Location"`
	VersionID         string    `json:"VersionID"`
	Expiration        time.Time `json:"Expiration"`
	ExpirationRuleId  string    `json:"ExpirationRuleId"`
	ChecksumCRC32     string    `json:"ChecksumCRC32"`
	ChecksumCRC32C    string    `json:"ChecksumCRC32C"`
	ChecksumSHA1      string    `json:"ChecksumSHA1"`
	ChecksumSHA256    string    `json:"ChecksumSHA256"`
	ChecksumCRC64NVME string    `json:"ChecksumCRC64NVME"`
	ChecksumMode      string    `json:"ChecksumMode"`
}

func NewObject() Object {
	return Object{}
}

func (o *Object) Set_LocalPath(path string) {
	const fileExtensionCSV string = ".csv"
	base := filepath.Base(path)
	if strings.Contains(base[len(base)-len(fileExtensionCSV):], fileExtensionCSV) {
		o.LocalName = base[:len(base)-len(fileExtensionCSV)]
		o.Name = o.LocalName
		o.LocalPath = path
	}
}

func (o *Object) Set_Name(name string) {
	o.Name = name
}

func (o *Object) FilePutObject(ctx context.Context, bucket *Bucket) (err error) {
	if bucket.Info.Name == "" {
		err = fmt.Errorf("Object: BucketName is empty")
		return
	}
	info, err := c.mc.FPutObject(ctx, bucket.Info.Name, o.Name, o.LocalPath, minio.PutObjectOptions{})
	if err != nil {
		err = fmt.Errorf("client.FPutObject: %s", err.Error())
		return
	}
	o.UploadInfo = Object_UploadInfo{
		Bucket:            info.Bucket,
		Key:               info.Key,
		ETag:              info.ETag,
		Size:              info.Size,
		LastModified:      info.LastModified,
		Location:          info.Location,
		VersionID:         info.VersionID,
		Expiration:        info.Expiration,
		ExpirationRuleId:  info.ExpirationRuleID,
		ChecksumCRC32:     info.ChecksumCRC32,
		ChecksumCRC32C:    info.ChecksumCRC32C,
		ChecksumSHA1:      info.ChecksumSHA1,
		ChecksumSHA256:    info.ChecksumSHA256,
		ChecksumCRC64NVME: info.ChecksumCRC64NVME,
		ChecksumMode:      info.ChecksumMode,
	}
	return
}
