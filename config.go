package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/osutils"
)

func Setup(app core.App) {
	settings := app.Settings()
	// APP
	settings.Meta.AppName = os.Getenv("APP_NAME")
	settings.Meta.AppURL = os.Getenv("APP_URL")

	// Proxy
	headerStr := os.Getenv("TRUSTED_PROXY_HEADERS")
	if headerStr != "" {
		// 将字符串按逗号分割成切片
		settings.TrustedProxy.Headers = strings.Split(headerStr, ",")
	}
	useLeftmost, _ := strconv.ParseBool(os.Getenv("USE_LEFT_MOST_IP"))
	settings.TrustedProxy.UseLeftmostIP = useLeftmost

	// Rate Limits
	enableRateLimits, _ := strconv.ParseBool(os.Getenv("ENABLE_RATE_LIMITS"))
	settings.RateLimits.Enabled = enableRateLimits
	settings.RateLimits.Rules = []core.RateLimitRule{
		{
			Label:       "*:auth",
			Duration:    3,
			MaxRequests: 2,
		},
		{
			Label:       "*:create",
			Duration:    5,
			MaxRequests: 20,
		},
		{
			Label:       "/api/batch",
			Duration:    1,
			MaxRequests: 3,
		}, {
			Label:       "/api/",
			Duration:    10,
			MaxRequests: 300,
		},
	}

	// Batch API
	settings.Batch.Enabled, _ = strconv.ParseBool(os.Getenv("ENABLE_BATCH_API"))
	settings.Batch.MaxRequests = 50
	settings.Batch.Timeout = 3

	// Hide Controls
	settings.Meta.HideControls = !osutils.IsProbablyGoRun()

	// Mail
	settings.Meta.SenderName = os.Getenv("SENDER_NAME")
	settings.Meta.SenderAddress = os.Getenv("SENDER_ADDRESS")
	settings.SMTP.AuthMethod = os.Getenv("SMTP_AUTH_METHOD") // PLAIN or LOGIN
	settings.SMTP.Host = os.Getenv("SMTP_HOST")
	if val := os.Getenv("SMTP_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			settings.SMTP.Port = port
		}
	}
	settings.SMTP.TLS, _ = strconv.ParseBool(os.Getenv("ENABLE_SMTP_TLS"))
	settings.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	settings.SMTP.LocalName = os.Getenv("SMTP_LOCALNAME")

	// S3
	settings.S3.Enabled, _ = strconv.ParseBool(os.Getenv("ENABLE_S3"))
	settings.S3.Region = os.Getenv("S3_REGION")
	settings.S3.Bucket = os.Getenv("S3_BUCKET")
	settings.S3.Endpoint = os.Getenv("S3_ENDPOINT")
	settings.S3.AccessKey = os.Getenv("S3_ACCESS_KEY")
	settings.S3.Secret = os.Getenv("S3_SECRET_KEY")
	settings.S3.ForcePathStyle, _ = strconv.ParseBool(os.Getenv("S3_FORCE_PATH_STYLE"))

	// Backup
	if os.Getenv("BACKUP_CRON") != "" {
		settings.Backups.Cron = os.Getenv("BACKUP_CRON")
		if val := os.Getenv("BACKUP_CRON_MAX_KEEP"); val != "" {
			if intVal, err := strconv.Atoi(val); err == nil {
				settings.Backups.CronMaxKeep = intVal
			}
		}
		settings.Backups.S3.Enabled, _ = strconv.ParseBool(os.Getenv("BACKUP_ENABLE_S3"))
		settings.Backups.S3.Region = os.Getenv("BACKUP_S3_REGION")
		settings.Backups.S3.Bucket = os.Getenv("BACKUP_S3_BUCKET")
		settings.Backups.S3.Endpoint = os.Getenv("BACKUP_S3_ENDPOINT")
		settings.Backups.S3.AccessKey = os.Getenv("BACKUP_S3_ACCESS_KEY")
		settings.Backups.S3.Secret = os.Getenv("BACKUP_S3_SECRET_KEY")
		settings.Backups.S3.ForcePathStyle, _ = strconv.ParseBool(os.Getenv("BACKUP_S3_FORCE_PATH_STYLE"))

	}

	settings.Logs.MaxDays = 7
	settings.Logs.MinLevel = 0
	settings.Logs.LogIP = true
	settings.Logs.LogAuthId = true

	if err := app.Save(settings); err != nil {
		app.Logger().Error("Save Globe Config Failed:", "error(s)", err)
	}
}
