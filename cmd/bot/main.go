package main

import (
	"log"
	"time"

	"decodeBot/internal/bot"
	"decodeBot/internal/client"
	"decodeBot/internal/config"
	"decodeBot/internal/scheduler"
	"decodeBot/internal/webhook"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize server client
	serverClient := client.NewServerClient(cfg.ServerURL, cfg.BotSecret)

	// Wait for server to be healthy before proceeding
	log.Println("⏳ Waiting for server to be ready...")
	maxRetries := 10
	retryDelay := 1 * time.Second
	serverReady := false

	for i := 0; i < maxRetries; i++ {
		if err := serverClient.HealthCheck(); err != nil {
			if i < maxRetries-1 {
				log.Printf("⚠️  Server not ready yet (attempt %d/%d), retrying in %v...", i+1, maxRetries, retryDelay)
				time.Sleep(retryDelay)
				// Exponential backoff, max 5 seconds
				retryDelay = retryDelay * 2
				if retryDelay > 5*time.Second {
					retryDelay = 5 * time.Second
				}
				continue
			}
			log.Printf("❌ Server health check failed after %d attempts: %v", maxRetries, err)
			log.Println("Bot will continue but server integration may not work")
		} else {
			serverReady = true
			log.Println("✓ Server connection established")
			break
		}
	}

	// Initialize bot
	pref := tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Debug {
		log.Println("🐛 Debug mode enabled")
	}

	log.Printf("✓ Bot authorized as @%s", b.Me.Username)

	// Initialize handler
	handler := bot.NewHandler(b, serverClient)

	// Register command handlers
	b.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			log.Printf("📥 Received update: %d | Text: %s", c.Update().ID, c.Text())
			return next(c)
		}
	})

	b.Handle("/start", handler.HandleStart)
	// Admin middleware
	adminOnly := func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Sender().ID != cfg.AdminID {
				// Silently ignore or reply? Users prefer silent ignore usually
				return nil
			}
			return next(c)
		}
	}

	// Admin commands
	b.Handle("/test_daily", handler.HandleTestDaily, adminOnly)
	b.Handle("/test_streak", handler.HandleTestStreak, adminOnly)
	b.Handle("/debug_schedule", handler.HandleDebugSchedule, adminOnly)

	// Initialize and start scheduler for daily notifications
	sched := scheduler.NewScheduler(b, serverClient)
	sched.Start()

	// Initialize and start webhook server for backend notifications
	webhookServer := webhook.NewServer(b, cfg.WebhookPort)
	webhookServer.Start()
	log.Printf("✓ Webhook server started on port %s", cfg.WebhookPort)

	// Send startup notification to admin only if server is ready
	if serverReady {
		bot.SendStartupNotification(b, serverClient, cfg.AdminID)
	} else {
		log.Println("⚠️  Skipping startup notification due to server connection issues")
	}

	log.Println("🤖 Bot is running...")
	log.Println("Press Ctrl+C to stop")

	// Start bot
	b.Start()
}
