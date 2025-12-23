package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	BotToken    string
	BotUsername string
	BotSecret   string
	ServerURL   string
	MiniAppURL  string
	Debug       bool
	AdminID     int64
	WebhookPort string // Port for webhook HTTP server
}

func Load() *Config {
	debug := os.Getenv("DEBUG") == "true"

	// Parse admin ID
	adminIDStr := os.Getenv("BOT_ADMIN_ID")
	var adminID int64
	if adminIDStr != "" {
		fmt.Sscanf(adminIDStr, "%d", &adminID)
	} else {
		adminID = 41361615 // Default fallback
	}

	webhookPort := os.Getenv("WEBHOOK_PORT")
	if webhookPort == "" {
		webhookPort = "8082" // Default webhook port
	}

	cfg := &Config{
		BotToken:    os.Getenv("BOT_TOKEN"),
		BotUsername: os.Getenv("BOT_USERNAME"),
		BotSecret:   os.Getenv("BOT_SECRET"),
		ServerURL:   os.Getenv("SERVER_URL"),
		MiniAppURL:  os.Getenv("MINI_APP_URL"),
		Debug:       debug,
		AdminID:     adminID,
		WebhookPort: webhookPort,
	}

	if cfg.BotToken == "" {
		log.Fatal("BOT_TOKEN is required")
	}

	if cfg.ServerURL == "" {
		cfg.ServerURL = "http://localhost:8081"
	}

	if cfg.MiniAppURL == "" {
		cfg.MiniAppURL = "https://ushpuras.dev/DEC0D3/"
	}

	return cfg
}
