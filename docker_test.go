package s3go_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var ports = []string{"9000", "9001"}

func StartDockerS3Server(t testing.TB) {
	t.Helper()
	ctx := context.Background()

	os.Setenv("RUSTFS_REGION", "eu-west-1")
	os.Setenv("RUSTFS_ACCESS_KEY_ID", "rustfsadmin")
	os.Setenv("RUSTFS_SECRET_ACCESS_KEY", "rustfsadmin")
	os.Setenv("RUSTFS_ENDPOINT_URL", "http://localhost:9000")

	req := testcontainers.ContainerRequest{
		Image: "rustfs/rustfs:latest",
		ExposedPorts: []string{
			fmt.Sprintf("%s:%s", ports[0], ports[0]),
			fmt.Sprintf("%s:%s", ports[1], ports[1]),
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(ports[0])).WithStartupTimeout(30*time.Second),
			wait.ForListeningPort(nat.Port(ports[1])).WithStartupTimeout(30*time.Second),
		),
		Env: map[string]string{
			"RUSTFS_ACCESS_KEY_ID":     "rustfsadmin",
			"RUSTFS_SECRET_ACCESS_KEY": "rustfsadmin",
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           log.New(os.Stdout, "testcontainers: ", log.LstdFlags),
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = container.Terminate(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})
}
