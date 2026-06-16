package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all application configuration, populated from environment variables.
type Config struct {
	App        AppConfig
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	AWS        AWSConfig
	SMTP       SMTPConfig
	Firebase   FirebaseConfig
	OpenSearch OpenSearchConfig
}

type AppConfig struct {
	Environment        string
	LogLevel           string
	LogFormat          string
	Debug              bool
	DefaultPageSize    int
	CORSAllowedOrigins []string
	DataEncryptionKey string
}

type FirebaseConfig struct {
	CredentialsPath string
}

type OpenSearchConfig struct {
	Endpoint string
	User     string
	Password string
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
	MaxConns int32
	MinConns int32
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	SecretKey     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

type AWSConfig struct {
	Region                   string
	AccessKeyID              string
	SecretAccessKey          string
	S3AssetsBucket           string // images, encrypted bundles (digimap, top_location, location_memo, latest-bundle)
	S3SnapshotBucket         string // snapshot draft/publish JSON blobs
	S3SyncBucket             string // venue sync data (VenueSyncFlag / SyncData JSON)
	CloudFrontDistributionID string
}

type SMTPConfig struct {
	Host      string
	Port      int
	User      string
	Password  string
	FromEmail string
	UseTLS    bool
}

// DSN builds the PostgreSQL connection string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&search_path=digimap_db,public",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

func (a AppConfig) IsProduction() bool { return a.Environment == "production" }
func (a AppConfig) IsDevelop() bool    { return a.Environment == "develop" }
func (a AppConfig) IsLocal() bool      { return a.Environment == "local" }

// Load reads configuration from environment variables (and optionally a .env file).
// Environment-specific .env files are loaded in order: .env.<environment>, then .env.
func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load .env file(s) — errors are intentionally ignored so missing files don't block startup
	env := v.GetString("APP_ENV")
	if env == "" {
		env = "local"
	}
	_ = godotenv.Load(".env." + env)
	_ = godotenv.Load(".env")

	// App defaults
	v.SetDefault("APP_ENV", "local")
	v.SetDefault("LOG_LEVEL", "debug")
	v.SetDefault("LOG_FORMAT", "text")
	v.SetDefault("DEBUG", false)
	v.SetDefault("DEFAULT_PAGE_SIZE", 20)
	v.SetDefault("CORS_ALLOWED_ORIGINS", "*")

	// Server defaults
	v.SetDefault("PORT", 8080)
	v.SetDefault("READ_TIMEOUT", "30s")
	v.SetDefault("WRITE_TIMEOUT", "30s")
	v.SetDefault("IDLE_TIMEOUT", "120s")

	// Database defaults
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_NAME", "digimap")
	v.SetDefault("DB_USER", "digimap")
	v.SetDefault("DB_PASSWORD", "")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_MAX_CONNS", 25)
	v.SetDefault("DB_MIN_CONNS", 5)

	// Redis defaults
	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	// JWT defaults
	v.SetDefault("JWT_SECRET_KEY", "")
	v.SetDefault("JWT_ACCESS_EXPIRY", "60m")
	v.SetDefault("JWT_REFRESH_EXPIRY", "24h")
	v.SetDefault("JWT_ISSUER", "digimap-backend")

	cfg := &Config{
		App: AppConfig{
			Environment:         v.GetString("APP_ENV"),
			LogLevel:            v.GetString("LOG_LEVEL"),
			LogFormat:           v.GetString("LOG_FORMAT"),
			Debug:               v.GetBool("DEBUG"),
			DefaultPageSize:     v.GetInt("DEFAULT_PAGE_SIZE"),
			CORSAllowedOrigins:  strings.Split(v.GetString("CORS_ALLOWED_ORIGINS"), ","),
			DataEncryptionKey: v.GetString("DATA_ENCRYPTION_KEY"),
		},
		Server: ServerConfig{
			Port:         v.GetInt("PORT"),
			ReadTimeout:  v.GetDuration("READ_TIMEOUT"),
			WriteTimeout: v.GetDuration("WRITE_TIMEOUT"),
			IdleTimeout:  v.GetDuration("IDLE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetInt("DB_PORT"),
			Name:     v.GetString("DB_NAME"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASSWORD"),
			SSLMode:  v.GetString("DB_SSLMODE"),
			MaxConns: int32(v.GetInt("DB_MAX_CONNS")),
			MinConns: int32(v.GetInt("DB_MIN_CONNS")),
		},
		Redis: RedisConfig{
			Addr:     v.GetString("REDIS_ADDR"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			SecretKey:     v.GetString("JWT_SECRET_KEY"),
			AccessExpiry:  v.GetDuration("JWT_ACCESS_EXPIRY"),
			RefreshExpiry: v.GetDuration("JWT_REFRESH_EXPIRY"),
			Issuer:        v.GetString("JWT_ISSUER"),
		},
		AWS: AWSConfig{
			Region:                   v.GetString("AWS_REGION"),
			AccessKeyID:              v.GetString("AWS_ACCESS_KEY_ID"),
			SecretAccessKey:          v.GetString("AWS_SECRET_ACCESS_KEY"),
			S3AssetsBucket:           v.GetString("AWS_S3_ASSETS_BUCKET"),
			S3SnapshotBucket:         v.GetString("AWS_S3_SNAPSHOT_BUCKET"),
			S3SyncBucket:             v.GetString("AWS_S3_SYNC_BUCKET"),
			CloudFrontDistributionID: v.GetString("AWS_CF_DISTRIBUTION_ID"),
		},
		SMTP: SMTPConfig{
			Host:      v.GetString("EMAIL_HOST"),
			Port:      v.GetInt("EMAIL_PORT"),
			User:      v.GetString("EMAIL_HOST_USER"),
			Password:  v.GetString("EMAIL_HOST_PASSWORD"),
			FromEmail: v.GetString("EMAIL_FROM_ADDR"),
			UseTLS:    v.GetBool("EMAIL_USE_TLS"),
		},
		Firebase: FirebaseConfig{
			CredentialsPath: v.GetString("FIREBASE_CREDENTIALS_PATH"),
		},
		OpenSearch: OpenSearchConfig{
			Endpoint: v.GetString("OPENSEARCH_ENDPOINT"),
			User:     v.GetString("OPENSEARCH_USER"),
			Password: v.GetString("OPENSEARCH_PASSWORD"),
		},
	}

	return cfg, nil
}
