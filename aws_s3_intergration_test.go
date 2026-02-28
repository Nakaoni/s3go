package s3go_test

import (
	"context"
	"s3go"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func TestAwsS3Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	StartDockerS3Server(t)

	ctx := context.Background()
	awsClient, _ := s3go.NewAwsS3(
		aws.Config{
			Region:       S3TestEnvs["RUSTFS_REGION"],
			BaseEndpoint: aws.String(S3TestEnvs["RUSTFS_ENDPOINT_URL"]),
			Credentials: aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(S3TestEnvs["RUSTFS_ACCESS_KEY_ID"], S3TestEnvs["RUSTFS_SECRET_ACCESS_KEY"], ""),
			),
		},
		true,
	)
	s3Server := s3go.NewS3Server(awsClient)

	t.Run("Empty List buckets", func(t *testing.T) {
		buckets, _, _ := s3Server.ListBuckets(ctx, nil)

		if len(buckets) != 0 {
			t.Errorf("Expected no buckets, got %v", len(buckets))
		}
	})
	t.Run("Create bucket", func(t *testing.T) {
		want := "bucket-1"
		ok, _ := s3Server.CreateBucket(ctx, want)

		if !ok {
			t.Errorf("Could not create bucket %q", want)
		}

		buckets, _, _ := s3Server.ListBuckets(ctx, nil)
		if len(buckets) != 1 {
			t.Errorf("Expected 1 bucket, got %v", len(buckets))
		}
	})
}
