package scheduler

import (
	"log"

	"decodeBot/internal/bot"
	"decodeBot/internal/client"

	"github.com/robfig/cron/v3"
	tele "gopkg.in/telebot.v4"
)

type Scheduler struct {
	cron   *cron.Cron
	bot    *tele.Bot
	client *client.ServerClient
}

func NewScheduler(bot *tele.Bot, serverClient *client.ServerClient) *Scheduler {
	return &Scheduler{
		cron:   cron.New(),
		bot:    bot,
		client: serverClient,
	}
}

// Start begins the scheduler
func (s *Scheduler) Start() {
	// Processing Queue every 2 minutes
	s.cron.AddFunc("*/2 * * * *", s.ProcessNotifications)

	// Trigger scheduling generation every hour (just to be safe and catch up)
	s.cron.AddFunc("0 * * * *", func() {
		if err := s.client.ScheduleNotifications(); err != nil {
			log.Printf("[SCHEDULER] Failed to trigger schedule: %v", err)
		}
	})

	s.cron.Start()
	log.Println("✓ Scheduler started - Smart Notification Queue enabled")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
}

// ProcessNotifications fetches and sends pending notifications
func (s *Scheduler) ProcessNotifications() {
	jobs, err := s.client.GetPendingNotifications(20) // Batch size 20
	if err != nil {
		log.Printf("[SCHEDULER] Failed to get jobs: %v", err)
		return
	}

	if len(jobs) == 0 {
		return
	}

	log.Printf("[SCHEDULER] Processing %d notification jobs...", len(jobs))

	for _, job := range jobs {
		if job.User == nil {
			log.Printf("[SCHEDULER] Job %d has no user data, skipping", job.ID)
			s.client.UpdateJobStatus(job.ID, "FAILED")
			continue
		}

		var message string
		if job.Type == "DAILY_CHALLENGE" {
			// Calculate best streak
			streak := 0
			if job.User.AllStreak > streak {
				streak = job.User.AllStreak
			}
			message = bot.GetDailyReminderMessage(job.User.FirstName, streak)
		} else {
			// Default fallback
			message = bot.GetDailyReminderMessage(job.User.FirstName, 0)
		}

		menu := bot.GetMainMenu()
		recipient := &tele.User{ID: job.User.TelegramID}

		if _, err := s.bot.Send(recipient, message, menu); err != nil {
			log.Printf("[SCHEDULER] Failed to send to %d: %v", job.User.TelegramID, err)

			// If blocked, maybe mark as FAILED or BLOCKED?
			// For now, marked as FAILED so we don't retry immediately (logic in server GetPending checks status=PENDING)
			s.client.UpdateJobStatus(job.ID, "FAILED")
		} else {
			log.Printf("[NOTIF] Sent to %s (@%s)", job.User.FirstName, job.User.Username)
			s.client.UpdateJobStatus(job.ID, "SENT")
		}
	}
}
