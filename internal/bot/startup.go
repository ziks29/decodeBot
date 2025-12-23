package bot

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"decodeBot/internal/client"

	tele "gopkg.in/telebot.v4"
)

var botStartTime = time.Now()

// StartupMetrics holds all the system information for the startup message
type StartupMetrics struct {
	ServerID         string
	StartTime        time.Time
	Uptime           time.Duration
	LastUpdate       time.Time
	WorkingDirectory string
	GoVersion        string
	Goroutines       int
	MemoryUsageMB    float64
	DBConnected      bool
	DBStatus         string
	TotalUsers       int
	ActiveUsers7d    int
}

// CollectStartupMetrics gathers all system metrics
func CollectStartupMetrics(serverClient *client.ServerClient) (*StartupMetrics, error) {
	metrics := &StartupMetrics{
		StartTime:  botStartTime,
		Uptime:     time.Since(botStartTime),
		GoVersion:  runtime.Version(),
		Goroutines: runtime.NumGoroutine(),
	}

	// Get hostname (server ID)
	if hostname, err := os.Hostname(); err == nil {
		metrics.ServerID = hostname
	} else {
		metrics.ServerID = "unknown"
	}

	// Get working directory
	if wd, err := os.Getwd(); err == nil {
		metrics.WorkingDirectory = wd
	} else {
		metrics.WorkingDirectory = "unknown"
	}

	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metrics.MemoryUsageMB = float64(m.Alloc) / 1024 / 1024

	// Check database connection
	if err := serverClient.HealthCheck(); err != nil {
		metrics.DBConnected = false
		metrics.DBStatus = "Disconnected"
	} else {
		metrics.DBConnected = true
		metrics.DBStatus = "Healthy"
	}

	// Get user statistics
	if stats, err := serverClient.GetUserStats(); err == nil {
		metrics.TotalUsers = stats.TotalUsers
		metrics.ActiveUsers7d = stats.ActiveUsers7d
	} else {
		log.Printf("Failed to get user stats: %v", err)
		metrics.TotalUsers = 0
		metrics.ActiveUsers7d = 0
	}

	// For last update, we'll use current time as a placeholder
	// In production, this could be set from build time or git commit timestamp
	metrics.LastUpdate = time.Now()

	return metrics, nil
}

// FormatStartupMessage creates a nicely formatted startup message
func FormatStartupMessage(metrics *StartupMetrics) string {
	// Format times
	startTimeStr := metrics.StartTime.Format("2006-01-02 15:04:05")
	lastUpdateStr := metrics.LastUpdate.Format("2006-01-02 15:04:05")
	
	// Format uptime
	uptimeStr := formatUptime(metrics.Uptime)

	// Database status emoji
	dbEmoji := "✅"
	if !metrics.DBConnected {
		dbEmoji = "❌"
	}

	return fmt.Sprintf(
		"🤖 Bot is active!\n\n"+
			"🖥️ Server: %s\n"+
			"⏰ Start time: %s\n"+
			"⌛️ Uptime: %s\n"+
			"🔄 Last update: %s\n"+
			"📂 Directory: %s\n\n"+
			"🔄 Go version: %s\n"+
			"⚙️ Goroutines: %d\n\n"+
			"💾 Memory usage: %.2f MB\n"+
			"🗄️ Database: %s\n"+
			"%s DB Status: %s\n\n"+
			"👥 Total users: %d\n"+
			"👤 Active users (7d): %d",
		metrics.ServerID,
		startTimeStr,
		uptimeStr,
		lastUpdateStr,
		metrics.WorkingDirectory,
		metrics.GoVersion,
		metrics.Goroutines,
		metrics.MemoryUsageMB,
		boolToStatus(metrics.DBConnected),
		dbEmoji,
		metrics.DBStatus,
		metrics.TotalUsers,
		metrics.ActiveUsers7d,
	)
}

// formatUptime converts duration to a human-readable format
func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// boolToStatus converts boolean to Connected/Disconnected
func boolToStatus(b bool) string {
	if b {
		return "Connected"
	}
	return "Disconnected"
}

// SendStartupNotification sends the startup message to the admin
func SendStartupNotification(bot *tele.Bot, serverClient *client.ServerClient, adminID int64) {
	if adminID == 0 {
		log.Println("⚠️  No admin ID configured, skipping startup notification")
		return
	}

	metrics, err := CollectStartupMetrics(serverClient)
	if err != nil {
		log.Printf("⚠️  Failed to collect startup metrics: %v", err)
		return
	}

	message := FormatStartupMessage(metrics)
	
	// Create recipient
	recipient := &tele.User{ID: adminID}
	
	if _, err := bot.Send(recipient, message); err != nil {
		log.Printf("⚠️  Failed to send startup notification to admin: %v", err)
	} else {
		log.Printf("✓ Startup notification sent to admin (ID: %d)", adminID)
	}
}
