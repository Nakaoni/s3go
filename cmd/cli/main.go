package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"s3go"
)

const (
	BIN           = "s3go"
	BUCKET_LIST   = "bucket:ls"
	OBJECT_LIST   = "object:ls"
	OBJECT_FIND   = "object:find"
	OBJECT_REMOVE = "object:rm"
)

var (
	commandsName = map[string]string{
		BUCKET_LIST:   "List buckets",
		OBJECT_LIST:   "List objects in a bucket",
		OBJECT_FIND:   "Find an object in a bucket",
		OBJECT_REMOVE: "Delete an object in a bucket",
	}
)

var client *s3go.Server

func main() {
	handle(os.Args[1:])
}

func handle(inputs []string) {
	if len(inputs) == 0 {
		showUsage()
		return
	}

	ctx := context.Background()
	c, err := getClientFromConfig(ctx)
	if err != nil {
		log.Println("Error while trying to instanciate s3editor")
		return
	}
	client = c

	command := inputs[0]
	args := inputs[1:]

	switch command {
	case BUCKET_LIST:
		listBuckets(ctx)
	case OBJECT_LIST:
		listObjects(ctx, args[0])
	case OBJECT_FIND:
		findObject(ctx, args[0], args[1])
	case OBJECT_REMOVE:
		removeObject(ctx, args[0], args[1])
	default:
		showUsage()
	}
}

func getClientFromConfig(ctx context.Context) (*s3go.Server, error) {
	awsS3, err := s3go.NewAwsS3(ctx)
	if err != nil {
		return nil, err
	}
	return s3go.NewS3Server(awsS3), nil
}

func showUsage() {
	fmt.Printf("Usage: %v [command] [<optional-arguments>]\n", BIN)
	fmt.Printf(" %v -l -- List available commands\n", BIN)
	fmt.Printf(" %v [command] -h -- List help of a command\n", BIN)

	fmt.Println("\nAvailable commands:")
	for name, description := range commandsName {
		fmt.Printf("\t%s -- %s\n", name, description)
	}
}

func listBuckets(ctx context.Context) {
	buckets, _, err := client.ListBuckets(ctx, nil)
	if err != nil {
		fmt.Println("Error while listing buckets:\n ", err)
	}

	fmt.Printf("Creation Date\tBucket Name\n")
	for _, b := range buckets {
		fmt.Printf("%v\t%v\n", b.CreationDate, *b.Name)
	}
}

func listObjects(ctx context.Context, bucketName string) {
	objects, _, err := client.ListObjects(ctx, bucketName, nil)
	if err != nil {
		fmt.Println("Error while listing objects of bucket: ", err)
	}
	fmt.Printf("\n\nLast Modified Date\tObject Name\n")
	for i, o := range objects {
		fmt.Printf("%v\t%v\n", o.LastModified, *o.Key)
		if i > 10 {
			break
		}
	}
}

func findObject(ctx context.Context, bucketName, objectKey string) {
	exists, err := client.ObjectExists(ctx, bucketName, objectKey)
	if err != nil {
		fmt.Printf("Error while searching object %q in bucket %q: \n", objectKey, bucketName)
		fmt.Println(err)
		return
	}
	if exists {
		fmt.Printf("Object %q exists in bucket %q\n", objectKey, bucketName)
	} else {
		fmt.Printf("Object %q does not exist in bucket %q\n", objectKey, bucketName)
	}
}

func removeObject(ctx context.Context, bucketName, objectKey string) {
	deleted, err := client.DeleteObject(ctx, bucketName, objectKey)
	if err != nil {
		fmt.Printf("Error while deleting object %q in bucket %q: \n", objectKey, bucketName)
		fmt.Println(err)
		return
	}
	if deleted {
		fmt.Printf("Object %q in bucket %q deleted\n", objectKey, bucketName)
	} else {
		fmt.Printf("Could not delete Object %q in bucket %q\n", objectKey, bucketName)
	}
}
