package s3go_test

import "testing"

func TestAwsS3Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	t.Run("List buckets", func(t *testing.T) {
		StartDockerS3Server(t)
	})
}
