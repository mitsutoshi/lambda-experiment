package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockS3Client struct{}

func (m *MockS3Client) GetObject(
	ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {

	if *params.Bucket == "error" {
		return nil, errors.New("TestError")
	} else if *params.Bucket == "timeout" {
		select {
		case <-time.After(2 * time.Second):
			return &s3.GetObjectOutput{}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return &s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewReader([]byte("mock data"))),
	}, nil
}

func (m *MockS3Client) PutObject(
	ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {

	if *params.Bucket == "error" {
		return nil, errors.New("TestError")
	} else if *params.Bucket == "timeout" {
		select {
		case <-time.After(2 * time.Second):
			return &s3.PutObjectOutput{}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return &s3.PutObjectOutput{}, nil
}

func (m *MockS3Client) HeadObject(
	ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	if *params.Bucket == "error" {
		return nil, errors.New("TestError")
	} else if *params.Bucket == "timeout" {
		select {
		case <-time.After(2 * time.Second):
			return &s3.HeadObjectOutput{}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return &s3.HeadObjectOutput{}, nil
}

func (m *MockS3Client) DeleteObject(
	ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	if *params.Bucket == "error" {
		return nil, errors.New("TestError")
	} else if *params.Bucket == "timeout" {
		select {
		case <-time.After(2 * time.Second):
			return &s3.DeleteObjectOutput{}, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return &s3.DeleteObjectOutput{}, nil
}

var (
	mock = &MockS3Client{}
	wait = 1 * time.Second
)

func TestMain(m *testing.M) {
	Client = mock
	code := m.Run()
	os.Exit(code)
}

func TestGetObject(t *testing.T) {

	testcases := []struct {
		bucket, key string
		wait        time.Duration
	}{
		{"normal", "testKey", wait},
		{"error", "testKey", wait},
		{"timeout", "testKey", wait},
	}

	for _, c := range testcases {

		// test
		obj, err := GetObject(c.bucket, c.key, c.wait)

		// check
		switch c.bucket {
		case "timeout":
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("expected: %v, actual: %s", context.DeadlineExceeded, err)
			}
		case "error":
			if err != nil && err.Error() != "TestError" {
				fmt.Println("2")
				t.Fatalf("expected: %s, actual: %s", "TestError", err)
			}
		case "normal":
			if obj == nil {
				t.Fatalf("obj is nil.")
			}
		}
	}
}

func TestPutObject(t *testing.T) {

	testcases := []struct {
		bucket, key, data, contentType string
		wait                           time.Duration
	}{
		{"normal", "testKey", "data", "image/jpeg", wait},
		{"error", "testKey", "data", "image/jpeg", wait},
		{"timeout", "testKey", "data", "image/jpeg", wait},
	}

	for _, c := range testcases {

		// test
		err := PutObject(c.bucket, c.key, strings.NewReader(c.data), c.contentType, c.wait)

		// check
		switch c.bucket {
		case "timeout":
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("expected: %v, actual: %s", context.DeadlineExceeded, err)
			}
		case "error":
			if err != nil && err.Error() != "TestError" {
				fmt.Println("2")
				t.Fatalf("expected: %s, actual: %s", "TestError", err)
			}
		}
	}
}

func TestHeadObject(t *testing.T) {

	testcases := []struct {
		bucket, key string
		wait        time.Duration
	}{
		{"normal", "testKey", wait},
		{"error", "testKey", wait},
		{"timeout", "testKey", wait},
	}

	for _, c := range testcases {

		// test
		obj, err := HeadObject(c.bucket, c.key, c.wait)

		// check
		switch c.bucket {
		case "timeout":
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("expected: %v, actual: %s", context.DeadlineExceeded, err)
			}
		case "error":
			if err != nil && err.Error() != "TestError" {
				fmt.Println("2")
				t.Fatalf("expected: %s, actual: %s", "TestError", err)
			}
		case "normal":
			if obj == nil {
				t.Fatalf("obj is nil.")
			}
		}
	}
}

func TestDeleteObject(t *testing.T) {

	testcases := []struct {
		bucket, key string
		wait        time.Duration
	}{
		{"normal", "testKey", wait},
		{"error", "testKey", wait},
		{"timeout", "testKey", wait},
	}

	for _, c := range testcases {

		// test
		obj, err := DeleteObject(c.bucket, c.key, c.wait)

		// check
		switch c.bucket {
		case "timeout":
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("expected: %v, actual: %s", context.DeadlineExceeded, err)
			}
		case "error":
			if err != nil && err.Error() != "TestError" {
				fmt.Println("2")
				t.Fatalf("expected: %s, actual: %s", "TestError", err)
			}
		case "normal":
			if obj == nil {
				t.Fatalf("obj is nil.")
			}
		}
	}
}
