package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func hPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// func hGetBucketInfo(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	m, err := minio.New().Init(ctx,
// 		conf.MinioEndpoint,
// 		conf.MinioRootUser,
// 		conf.MinioRootPassword,
// 		conf.BucketName,
// 		conf.MinioUseSSL,
// 	)
// 	if err != nil {
// 		log.Printf("hGetBucketInfo. Error:%s", err.Error())
// 		c.JSON(http.StatusInternalServerError, rError(err))
// 	}

// 	buckets, err := m.ListBucket(ctx)
// 	if err != nil {
// 		log.Printf("hGetBucketInfo. Error:%s", err.Error())
// 		c.JSON(http.StatusInternalServerError, rError(err))
// 	}

// 	c.JSON(http.StatusOK, buckets)
// }

// func hGetObjectsInfo(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	m, err := minio.New().Init(ctx,
// 		conf.MinioEndpoint,
// 		conf.MinioRootUser,
// 		conf.MinioRootPassword,
// 		conf.BucketName,
// 		conf.MinioUseSSL,
// 	)
// 	if err != nil {
// 		log.Printf("hGetObjectsInfo. Error:%s", err.Error())
// 		c.JSON(http.StatusInternalServerError, rError(err))
// 	}

// 	objects := m.ListObjects(ctx, conf.BucketName)
// 	c.JSON(http.StatusOK, objects)
// }

// func hCreateFile(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	fn := uuid.New().String()
// 	m, err := minio.New().Init(ctx,
// 		conf.MinioEndpoint,
// 		conf.MinioRootUser,
// 		conf.MinioRootPassword,
// 		conf.BucketName,
// 		conf.MinioUseSSL,
// 	)
// 	if err != nil {
// 		log.Printf("hGetObjectsInfo. Error:%s", err.Error())
// 		c.JSON(http.StatusInternalServerError, rError(err))
// 	}

// 	ct := c.Request.Header["Content-Type"][0]
// 	log.Printf("Content-Type:%s", ct)

// 	body, err := io.ReadAll(c.Request.Body)
// 	if err != nil {
// 		log.Printf("hCreateFile. Error:%s", err.Error())
// 		c.JSON(http.StatusInternalServerError, rError(err))
// 	}

// 	type info struct {
// 		Bucket            string    `json:"bucket"`
// 		Key               string    `json:"key"`
// 		ETag              string    `json:"etag"`
// 		Size              int64     `json:"size"`
// 		LastModified      time.Time `json:"last_modified"`
// 		Location          string    `json:"location"`
// 		VersionID         string    `json:"version_id"`
// 		Expiration        time.Time `json:"expiration"`
// 		ExpirationRuleID  string    `json:"expiration_rule_id"`
// 		ChecksumCRC32     string    `json:"ChecksumCRC32"`
// 		ChecksumCRC32C    string    `json:"ChecksumCRC32C"`
// 		ChecksumSHA1      string    `json:"ChecksumSHA1"`
// 		ChecksumSHA256    string    `json:"ChecksumSHA256"`
// 		ChecksumCRC64NVME string    `json:"ChecksumCRC64NVME"`
// 		ChecksumMode      string    `json:"ChecksumMode"`
// 	}

// 	size := int64(len(body))
// 	uploadInfo, url, err := m.CreateFile(ctx, conf.BucketName, fn, bytes.NewReader(body), size, ct)
// 	if err != nil {
// 		log.Printf("hCreateFile. Error:%s", err.Error())
// 		c.JSON(http.StatusInternalServerError, rError(err))
// 	}
// 	log.Printf("hCreateFile. URL:{%s} UploadFile:{%#v}", url.String(), uploadInfo)
// 	c.JSON(http.StatusOK, fmt.Appendf(nil, `"URL":"%s"`, url.String()))
// }
