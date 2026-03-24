package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/hhung06/digimap-backend/config"
)

// Storer generates pre-signed upload and download URLs for object storage.
type Storer interface {
	PresignUpload(ctx context.Context, key string, contentType string, ttl time.Duration) (string, error)
	PresignDownload(ctx context.Context, key string, ttl time.Duration) (string, error)
}

type s3Storer struct {
	client *s3.PresignClient
	bucket string
}

// NewS3Storer creates a Storer backed by AWS S3.
func NewS3Storer(cfg config.AWSConfig) (Storer, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	client := s3.NewFromConfig(awsCfg)
	presign := s3.NewPresignClient(client)
	return &s3Storer{client: presign, bucket: cfg.S3Bucket}, nil
}

func (s *s3Storer) PresignUpload(ctx context.Context, key, contentType string, ttl time.Duration) (string, error) {
	req, err := s.client.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("presign put: %w", err)
	}
	return req.URL, nil
}

func (s *s3Storer) PresignDownload(ctx context.Context, key string, ttl time.Duration) (string, error) {
	req, err := s.client.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("presign get: %w", err)
	}
	return req.URL, nil
}

// LogStorer is a no-op Storer for local development that returns fake URLs.
type LogStorer struct{}

func NewLogStorer() Storer { return &LogStorer{} }

func (s *LogStorer) PresignUpload(_ context.Context, key, _ string, _ time.Duration) (string, error) {
	return "https://s3.example.com/" + key + "?upload=1", nil
}

func (s *LogStorer) PresignDownload(_ context.Context, key string, _ time.Duration) (string, error) {
	return "https://s3.example.com/" + key, nil
}
