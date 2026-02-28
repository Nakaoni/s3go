package s3go_test

import (
	"context"
	"errors"
	"fmt"
	"s3go"
	"testing"
	"time"
)

type S3ClientStud struct {
	buckets []s3go.Bucket
	objects map[string][]s3go.Object
}

func (s *S3ClientStud) BucketsList(ctx context.Context, previousToken *string) (buckets []s3go.Bucket, continuationToken *string, err error) {
	buckets = s.buckets

	return
}

func (s *S3ClientStud) BucketAdd(ctx context.Context, name string) (bool, error) {
	s.buckets = append(s.buckets, s3go.Bucket{Name: &name})
	return true, nil
}

func (s *S3ClientStud) ObjectsList(ctx context.Context, bucketName string, previousToken *string) (objects []s3go.Object, continuationToken *string, err error) {
	objects = s.objects[bucketName]

	return
}

func (s *S3ClientStud) ObjectDelete(ctx context.Context, bucketName string, key string) (bool, error) {
	return true, nil
}

func (s *S3ClientStud) ObjectExists(ctx context.Context, bucketName string, key string) (bool, error) {
	objects, ok := s.objects[bucketName]

	if !ok {
		return false, errors.New("Does not exist")
	}

	for _, o := range objects {
		if *o.Key == key {
			return true, nil
		}
	}

	return false, nil
}

func TestBuckets(t *testing.T) {
	dummyS3Client := makeS3Client()
	server := s3go.NewS3Server(dummyS3Client)
	t.Run("List buckets", func(t *testing.T) {
		// TODO: test other returned values cases
		got, _, _ := server.ListBuckets(context.Background(), nil)

		want := dummyS3Client.buckets
		assertDeepEqual(t, got, want)
	})
	t.Run("Create bucket", func(t *testing.T) {
		bucketName := "bucket-1"
		ok, _ := server.CreateBucket(context.Background(), bucketName)

		if !ok {
			t.Errorf("Could not create bucket %q", bucketName)
		}

		got, _, _ := server.ListBuckets(context.Background(), nil)

		want := 11
		if len(got) != want {
			t.Errorf("Expected %d buckets, got %v", want, len(got))
		}
	})
}

func TestObjects(t *testing.T) {
	dummyS3Client := makeS3Client()
	server := s3go.NewS3Server(dummyS3Client)

	t.Run("Object list", func(t *testing.T) {
		bucketName := "bucket-1"
		// TODO: test other returned values cases
		got, _, _ := server.ListObjects(context.Background(), bucketName, nil)

		want := dummyS3Client.objects[bucketName]
		assertDeepEqual(t, got, want)

		if len(got) != len(want) {
			t.Errorf("Count: got %v, want %v", got, want)
		}
	})
	t.Run("Object Exists", func(t *testing.T) {
		bucketName := "bucket-1"
		objectKey := "object-1"
		// TODO: test other returned values cases
		got, _ := server.ObjectExists(context.Background(), bucketName, objectKey)

		if !got {
			t.Errorf("Could not find object %q in bucket %q", objectKey, bucketName)
		}
	})
}

func makeS3Client() *S3ClientStud {
	dummyS3Client := S3ClientStud{}

	buckets := make([]s3go.Bucket, 10)
	objectList := make(map[string][]s3go.Object)
	for i := range 10 {
		index := i + 1
		now := time.Now()
		bucketName := fmt.Sprintf("bucket-%d", index)
		region := fmt.Sprintf("euwest-%d", index)
		buckets[i] = s3go.Bucket{Name: &bucketName, CreationDate: &now, BucketRegion: &region}

		objects := make([]s3go.Object, 10)
		for j := range 10 {
			objectName := fmt.Sprintf("object-%d", j)
			size := int64(300)
			objects[j] = s3go.Object{Key: &objectName, LastModified: &now, Size: &size}
		}
		objectList[bucketName] = objects
	}

	dummyS3Client.buckets = buckets
	dummyS3Client.objects = objectList

	return &dummyS3Client
}
