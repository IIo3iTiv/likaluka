package minio

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
)

type Buckets []Bucket

type Bucket struct {
	Info             Bucket_Info                    `json:"Info"`
	Versioning       Bucket_VersioningConfiguration `json:"VersioningConfiguration"`
	Policy           Bucket_Policy                  `json:"Policy"`
	ObjectLookConfig Bucket_ObjectLookConfig        `json:"ObjectLookConfig"`
}

type Bucket_Info struct {
	Name         string    `json:"BucketName"`
	CreationDate time.Time `json:"CreationDate"`
	BucketRegion string    `json:"BucketRegion"`
}

type Bucket_VersioningConfiguration struct {
	XMLNameSpace     string   `json:"XmlNameSpace"`
	XMLNameLocal     string   `json:"XmlNameLocal"`
	Status           string   `json:"Status"`
	MFADelete        string   `json:"MfaDelete"`
	ExcludedPrefixes []string `json:"ExcludePrefixes"`
	ExcludeFolders   bool     `json:"ExcludeFolders"`
	PurgeOnDelete    string   `json:"PurgeOnDelete"`
}

type Bucket_Policy struct {
	CurrentPolicy string `json:"CurrentPolicy"`
}

type Bucket_ObjectLookConfig struct {
	ObjectLoock  string `json:"ObjectLock"`
	Mode         string `json:"Mode"`
	Validaty     uint   `json:"Validaty"`
	ValidatyUnit string `json:"ValidatyUnit"`
}

func NewBucket(ctx context.Context, name string) (bucket Bucket) {
	bucket.Info.Name = name
	return
}

func (b *Bucket) BucketExistsAndMake(ctx context.Context) (err error) {
	exists, err := b.BucketExists(ctx)
	if err != nil {
		return
	}
	if !exists {
		err = b.MakeBucket(ctx)
	}
	return
}

func (b *Bucket) MakeBucket(ctx context.Context) (err error) {
	return c.mc.MakeBucket(ctx, b.Info.Name, minio.MakeBucketOptions{})
}

func (b *Bucket) BucketExists(ctx context.Context) (exists bool, err error) {
	return c.mc.BucketExists(ctx, b.Info.Name)
}

func (b *Bucket) GetBucketVersioning(ctx context.Context) (err error) {
	v, err := c.mc.GetBucketVersioning(ctx, b.Info.Name)
	if err != nil {
		return
	}
	b.Versioning = Bucket_VersioningConfiguration{
		XMLNameSpace:   v.XMLName.Space,
		XMLNameLocal:   v.XMLName.Local,
		Status:         v.Status,
		MFADelete:      v.MFADelete,
		ExcludeFolders: v.ExcludeFolders,
		PurgeOnDelete:  v.PurgeOnDelete,
	}
	for _, p := range v.ExcludedPrefixes {
		b.Versioning.ExcludedPrefixes = append(b.Versioning.ExcludedPrefixes, p.Prefix)
	}
	return
}

func (b *Bucket) GetBucketPolicy(ctx context.Context) (err error) {
	cp, err := c.mc.GetBucketPolicy(ctx, b.Info.Name)
	if err != nil {
		return
	}
	b.Policy.CurrentPolicy = cp
	return
}

func (b *Bucket) GetObjectLookConfig(ctx context.Context) (err error) {
	ol, m, v, u, err := c.mc.GetObjectLockConfig(ctx, b.Info.Name)
	if err != nil {
		return err
	}
	b.ObjectLookConfig = Bucket_ObjectLookConfig{
		ObjectLoock:  ol,
		Mode:         m.String(),
		Validaty:     *v,
		ValidatyUnit: u.String(),
	}
	return
}

func GetListBuckets(ctx context.Context) (buckets Buckets, err error) {
	lb, err := c.mc.ListBuckets(ctx)
	for _, b := range lb {
		buckets = append(buckets, Bucket{
			Info: Bucket_Info{
				Name:         b.Name,
				CreationDate: b.CreationDate,
				BucketRegion: b.BucketRegion,
			},
		})
	}
	return
}

func (buckets Buckets) Prepare(ctx context.Context) (err error) {
	for i := range buckets {
		err = nil
		err = fmt.Errorf("GetBucketVersioning:{%s} ", buckets[i].GetBucketVersioning(ctx).Error())
		err = fmt.Errorf("GetBucketPolicy:{%s} ", buckets[i].GetBucketPolicy(ctx).Error())
		err = fmt.Errorf("GetBucketVersioning:{%s} ", buckets[i].GetBucketVersioning(ctx).Error())
		if err != nil {
			return
		}
	}
	return
}
