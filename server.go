package s3go

import (
	"context"
	"time"
)

type Bucket struct {
	BucketRegion *string
	CreationDate *time.Time
	Name         *string
}

type Object struct {
	Key          *string
	LastModified *time.Time
	Size         *int64
}

type S3Client interface {
	BucketsList(ctx context.Context, previousToken *string) (buckets []Bucket, continuationToken *string, err error)
	ObjectsList(ctx context.Context, bucketName string, previousToken *string) (objects []Object, continuationToken *string, err error)
	ObjectDelete(ctx context.Context, bucketName string, key string) (bool, error)
	ObjectExists(ctx context.Context, bucketName string, key string) (bool, error)
}

type Server struct {
	s3client S3Client
}

func NewS3Server(client S3Client) *Server {
	return &Server{client}
}

func (s Server) ListBuckets(ctx context.Context, previousToken *string) (buckets []Bucket, continuationToken *string, err error) {
	return s.s3client.BucketsList(ctx, previousToken)
}

func (s Server) ListObjects(ctx context.Context, bucketName string, previousToken *string) (objects []Object, continuationToken *string, err error) {
	return s.s3client.ObjectsList(ctx, bucketName, previousToken)
}

func (s Server) DeleteObject(ctx context.Context, bucketName string, key string) (bool, error) {
	return s.s3client.ObjectDelete(ctx, bucketName, key)
}

func (s Server) ObjectExists(ctx context.Context, bucketName string, key string) (bool, error) {
	return s.s3client.ObjectExists(ctx, bucketName, key)
}
