package s3

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3IF interface {
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

var Client S3IF

func HeadObject(bucket, key string, timeout time.Duration) (*s3.HeadObjectOutput, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	headObj, err := Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return headObj, nil
}

func GetObject(bucket, key string, timeout time.Duration) (*s3.GetObjectOutput, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	obj, err := Client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func PutObject(bucket, key string, r io.Reader, contentType string, timeout time.Duration) error {

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        r,
		ContentType: aws.String(contentType),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := Client.PutObject(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

func DeleteObject(bucket string, key string) error {

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := Client.DeleteObject(ctx, input)
	if err != nil {
		return err
	}
	return nil
}
