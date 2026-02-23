package s3go

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type AwsS3 struct {
	*s3.Client
}

func NewAwsS3(cfg aws.Config) (*AwsS3, error) {
	return &AwsS3{s3.NewFromConfig(cfg)}, nil
}

func (s *AwsS3) BucketsList(ctx context.Context, previousToken *string) (buckets []Bucket, continuationToken *string, err error) {
	var output *s3.ListBucketsOutput
	bucketPaginator := s3.NewListBucketsPaginator(s, &s3.ListBucketsInput{})
	if bucketPaginator.HasMorePages() {
		output, err = bucketPaginator.NextPage(ctx)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
				fmt.Println("You don't have permission to list buckets for this account.")
				err = apiErr
			} else {
				log.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
			}
		}
		for _, b := range output.Buckets {
			buckets = append(buckets, Bucket{BucketRegion: b.BucketRegion, Name: b.Name, CreationDate: b.CreationDate})
		}
	}
	return buckets, output.ContinuationToken, err
}

func (s *AwsS3) ObjectsList(ctx context.Context, bucketName string, previousToken *string) (objects []Object, continuationToken *string, err error) {
	var output *s3.ListObjectsV2Output

	objectPaginator := s3.NewListObjectsV2Paginator(s, &s3.ListObjectsV2Input{Bucket: aws.String(bucketName), ContinuationToken: previousToken})
	if objectPaginator.HasMorePages() {
		output, err = objectPaginator.NextPage(ctx)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
				fmt.Println("You don't have permission to list objects in this bucket for this account.")
				err = apiErr
			} else {
				log.Printf("Couldn't list objects in this bucket for your account. Here's why: %v\n", err)
			}
		}
		for _, o := range output.Contents {
			objects = append(objects, Object{Key: o.Key, LastModified: o.LastModified, Size: o.Size})
		}
		continuationToken = output.ContinuationToken
	}

	return
}

func (s *AwsS3) ObjectDelete(ctx context.Context, bucketName string, key string) (bool, error) {
	exists, err := s.ObjectExists(ctx, bucketName, key)
	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	_, err = s.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(bucketName), Key: aws.String(key)})
	if err != nil {
		log.Printf("Could not delete object. Here's why: %v\n", err)
	}

	return true, nil
}

func (s *AwsS3) ObjectExists(ctx context.Context, bucketName string, key string) (bool, error) {
	_, err := s.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(key)})
	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			fmt.Printf("Object with key %q in bucket %q does not exists\n", key, bucketName)
		} else {
			log.Printf("Could not find object. Here's why: %v\n", err)
		}
		return false, err
	}

	return true, nil
}
