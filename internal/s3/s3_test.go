package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockS3Client struct{}

func (m *MockS3Client) GetObject(
	ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if *params.Bucket == "error" {
		return nil, errors.New("GetObjectError")
	}
	data := "mock data"
	body := io.NopCloser(bytes.NewReader([]byte(data)))
	return &s3.GetObjectOutput{
		Body: body,
	}, nil
}

func (m *MockS3Client) PutObject(
	ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if *params.Bucket == "error" {
		return nil, errors.New("PutObjectError")
	}
	return &s3.PutObjectOutput{}, nil
}

func (m *MockS3Client) HeadObject(
	ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	if *params.Bucket == "error" {
		return nil, errors.New("HeadObjectError")
	}
	return &s3.HeadObjectOutput{}, nil
}

func (m *MockS3Client) DeleteObject(
	ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	if *params.Bucket == "error" {
		return nil, errors.New("DeleteObjectError")
	}
	return &s3.DeleteObjectOutput{}, nil
}

var mock = &MockS3Client{}

// TODO test timeout
func TestGetObjectFromS3(t *testing.T) {
	Client = mock

	testcases := []struct {
		bucket, key string
	}{
		{"normal", "testKey"},
		{"error", "testKey"},
	}

	for _, c := range testcases {
		obj, err := GetObjectFromS3(c.bucket, c.key)

		if err != nil {
			if err.Error() != "GetObjectError" {
				t.Fatalf("expected: %s, actual: %s", "GetObjectError", err)
			}
		} else {
			if obj == nil {
				t.Fatalf("obj is nil.")
			}
		}
	}
}

func TestPutObjectToS3(t *testing.T) {
	Client = mock

	testcases := []struct {
		bucket, key, data, contentType string
	}{
		{"normal", "testKey", "data", "image/jpeg"},
		{"error", "testKey", "data", "image/jpeg"},
	}

	for _, c := range testcases {
		r := strings.NewReader(c.data)
		err := PutObjectToS3(c.bucket, c.key, r, c.contentType)
		if err != nil {
			if err.Error() != "PutObjectError" {
				t.Fatalf("expected: %s, actual: %s", "GetObjectError", err)
			}
		}
	}
}
