package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/osutils"
)

// Helper Functions

func setStringIfPresent(key string, dst *string) {
	if v, ok := os.LookupEnv(key); ok {
		*dst = v
	}
}

func setBoolIfPresent(key string, dst *bool) {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			*dst = b
		}
	}
}

func setIntIfPresent(key string, dst *int) {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			*dst = i
		}
	}
}

// Setup
func Setup(app core.App) {
	settings := app.Settings()

	// APP
	setStringIfPresent("APP_NAME", &settings.Meta.AppName)
	setStringIfPresent("APP_URL", &settings.Meta.AppURL)

	// Rate Limits
	setBoolIfPresent("ENABLE_RATE_LIMITS", &settings.RateLimits.Enabled)
	settings.RateLimits.Rules = []core.RateLimitRule{
		{Label: "*:auth", Duration: 3, MaxRequests: 2},
		{Label: "*:create", Duration: 5, MaxRequests: 20},
		{Label: "/api/batch", Duration: 1, MaxRequests: 3},
		{Label: "/api/", Duration: 10, MaxRequests: 300},
	}

	// --- Proxy ---

	if v, ok := os.LookupEnv("TRUSTED_PROXY_HEADERS"); ok {
		headers := []string{}
		for _, h := range strings.Split(v, ",") {
			if trimmed := strings.TrimSpace(h); trimmed != "" {
				headers = append(headers, trimmed)
			}
		}
		settings.TrustedProxy.Headers = headers
	}
	setBoolIfPresent("USE_LEFT_MOST_IP", &settings.TrustedProxy.UseLeftmostIP)

	// Batch API
	setBoolIfPresent("ENABLE_BATCH_API", &settings.Batch.Enabled)
	setIntIfPresent("BATCH_MAX_REQUESTS", &settings.Batch.MaxRequests)

	settings.Meta.HideControls = !osutils.IsProbablyGoRun()

	// SMTP
	setStringIfPresent("SENDER_NAME", &settings.Meta.SenderName)
	setStringIfPresent("SENDER_ADDRESS", &settings.Meta.SenderAddress)
	setBoolIfPresent("ENABLE_SMTP", &settings.SMTP.Enabled)
	setStringIfPresent("SMTP_AUTH_METHOD", &settings.SMTP.AuthMethod)
	setStringIfPresent("SMTP_HOST", &settings.SMTP.Host)
	setIntIfPresent("SMTP_PORT", &settings.SMTP.Port)
	setBoolIfPresent("ENABLE_SMTP_TLS", &settings.SMTP.TLS)
	setStringIfPresent("SMTP_USERNAME", &settings.SMTP.Username)
	setStringIfPresent("SMTP_PASSWORD", &settings.SMTP.Password)
	setStringIfPresent("SMTP_LOCALNAME", &settings.SMTP.LocalName)

	// S3
	setBoolIfPresent("ENABLE_S3", &settings.S3.Enabled)
	setStringIfPresent("S3_REGION", &settings.S3.Region)
	setStringIfPresent("S3_BUCKET", &settings.S3.Bucket)
	setStringIfPresent("S3_ENDPOINT", &settings.S3.Endpoint)
	setStringIfPresent("S3_ACCESS_KEY", &settings.S3.AccessKey)
	setStringIfPresent("S3_SECRET_KEY", &settings.S3.Secret)
	setBoolIfPresent("S3_FORCE_PATH_STYLE", &settings.S3.ForcePathStyle)

	// Backup
	setStringIfPresent("BACKUP_CRON", &settings.Backups.Cron)
	setIntIfPresent("BACKUP_CRON_MAX_KEEP", &settings.Backups.CronMaxKeep)
	setBoolIfPresent("BACKUP_ENABLE_S3", &settings.Backups.S3.Enabled)
	setStringIfPresent("BACKUP_S3_REGION", &settings.Backups.S3.Region)
	setStringIfPresent("BACKUP_S3_BUCKET", &settings.Backups.S3.Bucket)
	setStringIfPresent("BACKUP_S3_ENDPOINT", &settings.Backups.S3.Endpoint)
	setStringIfPresent("BACKUP_S3_ACCESS_KEY", &settings.Backups.S3.AccessKey)
	setStringIfPresent("BACKUP_S3_SECRET_KEY", &settings.Backups.S3.Secret)
	setBoolIfPresent("BACKUP_S3_FORCE_PATH_STYLE", &settings.Backups.S3.ForcePathStyle)

	// Logs
	settings.Logs.MaxDays = 7
	settings.Logs.MinLevel = 0
	settings.Logs.LogIP = true
	settings.Logs.LogAuthId = true

	// Save
	if err := app.Save(settings); err != nil {
		app.Logger().Error("Save Global Config Failed:", "err", err)
	}
}
