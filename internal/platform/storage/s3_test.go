package storage

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/hhung06/digimap-backend/config"
)

type recordingPutObjectClient struct {
	input *s3.PutObjectInput
	body  string
}

func (c *recordingPutObjectClient) PutObject(_ context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	c.input = input
	body, err := io.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}
	c.body = string(body)
	return &s3.PutObjectOutput{}, nil
}

func TestS3StorerPutMediaForwardsObjectMetadata(t *testing.T) {
	client := &recordingPutObjectClient{}
	storer := &s3Storer{media: client, bucket: "media-bucket"}

	err := storer.PutMedia(context.Background(), "develop/media/articles/id/images/upload.png", "image/png", 7, strings.NewReader("content"))
	if err != nil {
		t.Fatalf("PutMedia returned error: %v", err)
	}
	if got := *client.input.Bucket; got != "media-bucket" {
		t.Fatalf("expected bucket media-bucket, got %q", got)
	}
	if got := *client.input.Key; got != "develop/media/articles/id/images/upload.png" {
		t.Fatalf("unexpected key %q", got)
	}
	if got := *client.input.ContentType; got != "image/png" {
		t.Fatalf("expected content type image/png, got %q", got)
	}
	if got := *client.input.ContentLength; got != 7 {
		t.Fatalf("expected content length 7, got %d", got)
	}
	if client.body != "content" {
		t.Fatalf("expected forwarded body content, got %q", client.body)
	}
}

func TestLoadAWSConfigUsesDefaultCredentialChain(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "environment-access-key")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "environment-secret-key")
	t.Setenv("AWS_SESSION_TOKEN", "environment-session-token")

	awsCfg, err := loadAWSConfig(context.Background(), config.AWSConfig{
		Region:          "ap-southeast-1",
		AccessKeyID:     "legacy-static-access-key",
		SecretAccessKey: "legacy-static-secret-key",
	})
	if err != nil {
		t.Fatalf("loadAWSConfig returned error: %v", err)
	}
	credentials, err := awsCfg.Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("retrieve credentials: %v", err)
	}
	if credentials.AccessKeyID != "environment-access-key" {
		t.Fatalf("expected default chain environment credentials, got %q", credentials.AccessKeyID)
	}
	if awsCfg.Region != "ap-southeast-1" {
		t.Fatalf("expected configured region, got %q", awsCfg.Region)
	}
}
