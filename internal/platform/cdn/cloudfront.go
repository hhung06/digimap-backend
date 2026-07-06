package cdn

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"

	"github.com/hhung06/digimap-backend/config"
)

// Invalidator sends CloudFront cache invalidations.
// Mirrors Django's invalidate_cloudfront (indoormap-backend/utils/s3services.py:85).
// Paths must include the leading slash (e.g., "/production/top_location/public/...").
type Invalidator interface {
	Invalidate(ctx context.Context, paths []string) (id string, err error)
}

type cfInvalidator struct {
	client         *cloudfront.Client
	distributionID string
}

// NewCloudFrontInvalidator creates an Invalidator backed by AWS CloudFront.
func NewCloudFrontInvalidator(cfg config.AWSConfig) (Invalidator, error) {
	if cfg.CloudFrontDistributionID == "" {
		return nil, fmt.Errorf("cloudfront distribution id is required (set AWS_CF_DISTRIBUTION_ID)")
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config for cloudfront: %w", err)
	}
	return &cfInvalidator{
		client:         cloudfront.NewFromConfig(awsCfg),
		distributionID: cfg.CloudFrontDistributionID,
	}, nil
}

// Invalidate creates a CloudFront invalidation for the given paths (fire-and-forget; non-fatal errors are returned
// but callers should log and continue). Mirrors Django's fire-and-forget pattern (publish_venue_v2.py:926-931).
func (c *cfInvalidator) Invalidate(ctx context.Context, paths []string) (string, error) {
	if len(paths) == 0 {
		return "", nil
	}
	// CloudFront caps at 3000 paths per call.
	if len(paths) > 3000 {
		paths = paths[:3000]
	}
	items := make([]string, len(paths))
	for i, p := range paths {
		items[i] = encodeInvalidationPath(p)
	}
	ref := strconv.FormatInt(time.Now().UnixNano(), 10)
	out, err := c.client.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(c.distributionID),
		InvalidationBatch: &cftypes.InvalidationBatch{
			Paths: &cftypes.Paths{
				Quantity: aws.Int32(int32(len(items))),
				Items:    items,
			},
			CallerReference: aws.String(ref),
		},
	})
	if err != nil {
		return "", fmt.Errorf("cloudfront invalidation: %w", err)
	}
	if out.Invalidation != nil && out.Invalidation.Id != nil {
		return *out.Invalidation.Id, nil
	}
	return "", nil
}

// encodeInvalidationPath URL-encodes each path segment (spaces, unicode, and other
// characters CloudFront rejects unencoded — e.g. from free-text theme names embedded
// in S3 keys) while preserving "/" as the literal path separator, per AWS's
// requirement to percent-encode non-ASCII/unsafe characters in invalidation paths.
func encodeInvalidationPath(p string) string {
	segments := strings.Split(p, "/")
	for i, seg := range segments {
		segments[i] = url.PathEscape(seg)
	}
	return strings.Join(segments, "/")
}
