package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hhung06/digimap-backend/config"
)

func TestLoadUsesLegacyS3BucketAsAssetsBucketFallback(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("AWS_S3_BUCKET", "legacy-assets-bucket")
	t.Setenv("AWS_S3_ASSETS_BUCKET", "")

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, "legacy-assets-bucket", cfg.AWS.S3AssetsBucket)
}

func TestLoadPrefersExplicitAssetsBucket(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("AWS_S3_BUCKET", "legacy-assets-bucket")
	t.Setenv("AWS_S3_ASSETS_BUCKET", "explicit-assets-bucket")

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, "explicit-assets-bucket", cfg.AWS.S3AssetsBucket)
}

func TestMain(m *testing.M) {
	_ = os.Unsetenv("AWS_S3_BUCKET")
	_ = os.Unsetenv("AWS_S3_ASSETS_BUCKET")
	os.Exit(m.Run())
}
