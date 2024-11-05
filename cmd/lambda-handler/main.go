package main

import (
	"context"
	"log"
	"resizeimage/internal/resize"
	s3if "resizeimage/internal/s3"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func init() {

	// set log conf
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// TODO
	// set s3 client
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}
	s3if.Client = s3.NewFromConfig(cfg)

}

func handleRequest(ctx context.Context, event events.S3Event) error {
	for _, e := range event.Records {
		log.Printf("EVENT => %v\n", event)
		err := resize.HandleS3Event(e)
		if err != nil {
			// TODO remained records will not be processed
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
