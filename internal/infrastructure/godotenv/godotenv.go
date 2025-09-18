package godotenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	// Server
	Environment string
	HTTPPort    string

	// Google Cloud / Vertex AI
	GoogleCloudProjectID   string
	GoogleApplicationCreds string
	GoogleCloudRegion      string
	GoogleVeoModel         string
	GoogleVeoMaxDuration   string
	GoogleAPIKey           string

	// Dataclay QUE
	DataclayQueHost           string
	DataclayQueAPIKey         string
	DataclayQueSatelliteID    string
	DataclayQueTemplaterBotID string
	DataclayQueAETemplate     string

	// AWS
	AWSAccessKeyID      string
	AWSSecretAccessKey  string
	AWSRegion           string
	AWSS3Bucket         string
	AWSS3StagingBucket  string
	AWSCFDistributionID string
	AWSCFKeyPairID      string
	AWSCFPrivateKeyPath string

	// Email / Postmark
	PostmarkAPIKey     string
	PostmarkFromEmail  string
	PostmarkFromName   string
	PostmarkTemplateID string
	SkipEmailSending   string

	// Database
	DatabaseURL string

	// Redis
	RedisURL string

	// Security & Rate Limiting
	RecaptchaSiteKey             string
	RecaptchaSecretKey           string
	RateLimitWindowMS            string
	RateLimitMaxRequests         string
	RateLimitMaxRequestsPerEmail string
	RateLimitWindowPerEmailMS    string

	// File Upload
	MaxFileSize      string
	AllowedFileTypes string
	UploadTimeout    string
	WatermarkPath    string

	// Video Processing
	VideoTargetWidth   string
	VideoTargetHeight  string
	VideoTargetBitrate string
	VideoTargetFPS     string
	VideoMaxDuration   string
	VideoMinDuration   string

	// Stripe
	StripeSecretKey      string
	StripePublishableKey string
	StripeWebhookSecret  string
	StripePriceID        string
}

func NewEnv() *Env {
	e := &Env{}
	if err := e.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
	return e
}

func (e *Env) Load() error {
	if err := godotenv.Load(".env"); err != nil {
		// Fallback to system env if .env is missing
		log.Printf("No .env file found, falling back to system environment")
	}

	// Server
	e.Environment = os.Getenv("ENVIRONMENT")
	e.HTTPPort = os.Getenv("HTTP_PORT")

	// Google Cloud
	e.GoogleCloudProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT_ID")
	e.GoogleApplicationCreds = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	e.GoogleCloudRegion = os.Getenv("GOOGLE_CLOUD_REGION")
	e.GoogleVeoModel = os.Getenv("GOOGLE_VEO_MODEL")
	e.GoogleVeoMaxDuration = os.Getenv("GOOGLE_VEO_MAX_DURATION")
	e.GoogleAPIKey = os.Getenv("GOOGLE_API_KEY")

	// Dataclay
	e.DataclayQueHost = os.Getenv("DATACLAY_QUE_HOST")
	e.DataclayQueAPIKey = os.Getenv("DATACLAY_QUE_API_KEY")
	e.DataclayQueSatelliteID = os.Getenv("DATACLAY_QUE_SATELLITE_ID")
	e.DataclayQueTemplaterBotID = os.Getenv("DATACLAY_QUE_TEMPLATER_BOT_ID")
	e.DataclayQueAETemplate = os.Getenv("DATACLAY_QUE_AE_TEMPLATE")

	// AWS
	e.AWSAccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	e.AWSSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	e.AWSRegion = os.Getenv("AWS_REGION")
	e.AWSS3Bucket = os.Getenv("AWS_S3_BUCKET")
	e.AWSS3StagingBucket = os.Getenv("AWS_S3_STAGING_BUCKET")
	e.AWSCFDistributionID = os.Getenv("AWS_CLOUDFRONT_DISTRIBUTION_ID")
	e.AWSCFKeyPairID = os.Getenv("AWS_CLOUDFRONT_KEY_PAIR_ID")
	e.AWSCFPrivateKeyPath = os.Getenv("AWS_CLOUDFRONT_PRIVATE_KEY_PATH")

	// Postmark
	e.PostmarkAPIKey = os.Getenv("POSTMARK_API_KEY")
	e.PostmarkFromEmail = os.Getenv("POSTMARK_FROM_EMAIL")
	e.PostmarkFromName = os.Getenv("POSTMARK_FROM_NAME")
	e.PostmarkTemplateID = os.Getenv("POSTMARK_TEMPLATE_ID")
	e.SkipEmailSending = os.Getenv("SKIP_EMAIL_SENDING")

	// Database
	e.DatabaseURL = os.Getenv("DATABASE_URL")

	// Redis
	e.RedisURL = os.Getenv("REDIS_URL")

	// Security & Rate Limiting
	e.RecaptchaSiteKey = os.Getenv("RECAPTCHA_SITE_KEY")
	e.RecaptchaSecretKey = os.Getenv("RECAPTCHA_SECRET_KEY")
	e.RateLimitWindowMS = os.Getenv("RATE_LIMIT_WINDOW_MS")
	e.RateLimitMaxRequests = os.Getenv("RATE_LIMIT_MAX_REQUESTS")
	e.RateLimitMaxRequestsPerEmail = os.Getenv("RATE_LIMIT_MAX_REQUESTS_PER_EMAIL")
	e.RateLimitWindowPerEmailMS = os.Getenv("RATE_LIMIT_WINDOW_PER_EMAIL_MS")

	// File Upload
	e.MaxFileSize = os.Getenv("MAX_FILE_SIZE")
	e.AllowedFileTypes = os.Getenv("ALLOWED_FILE_TYPES")
	e.UploadTimeout = os.Getenv("UPLOAD_TIMEOUT")
	e.WatermarkPath = os.Getenv("WATERMARK_PATH")

	// Video
	e.VideoTargetWidth = os.Getenv("VIDEO_TARGET_WIDTH")
	e.VideoTargetHeight = os.Getenv("VIDEO_TARGET_HEIGHT")
	e.VideoTargetBitrate = os.Getenv("VIDEO_TARGET_BITRATE")
	e.VideoTargetFPS = os.Getenv("VIDEO_TARGET_FPS")
	e.VideoMaxDuration = os.Getenv("VIDEO_MAX_DURATION")
	e.VideoMinDuration = os.Getenv("VIDEO_MIN_DURATION")

	// Stripe
	e.StripeSecretKey = os.Getenv("STRIPE_SECRET_KEY")
	e.StripePublishableKey = os.Getenv("STRIPE_PUBLISHABLE_KEY")
	e.StripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")
	e.StripePriceID = os.Getenv("STRIPE_PRICE_ID")

	return nil
}
