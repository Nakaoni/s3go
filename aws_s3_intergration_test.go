package s3go_test

import (
	"context"
	"os"
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

	var (
		region            = os.Getenv("RUSTFS_REGION")
		access_key_id     = os.Getenv("RUSTFS_ACCESS_KEY_ID")
		secret_access_key = os.Getenv("RUSTFS_SECRET_ACCESS_KEY")
		endpoint          = os.Getenv("RUSTFS_ENDPOINT_URL")
	)

	ctx := context.Background()
	awsClient, _ := s3go.NewAwsS3(aws.Config{
		Region:       region,
		BaseEndpoint: aws.String(endpoint),
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(access_key_id, secret_access_key, ""),
		),
	})
	s3Server := s3go.NewS3Server(awsClient)

	t.Run("List buckets", func(t *testing.T) {
		buckets, _, _ := s3Server.ListBuckets(ctx, nil)

		if len(buckets) != 0 {
			t.Errorf("Expected no buckets, got %v", len(buckets))
		}
	})
}
