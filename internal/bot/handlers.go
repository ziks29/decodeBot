package bot

import (
	"fmt"
	"log"

	"decodeBot/internal/client"
	"decodeBot/internal/models"

	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	bot    *tele.Bot
	client *client.ServerClient
}

func NewHandler(bot *tele.Bot, serverClient *client.ServerClient) *Handler {
	return &Handler{
		bot:    bot,
		client: serverClient,
	}
}

// HandleStart handles the /start command
func (h *Handler) HandleStart(c tele.Context) error {
	user := c.Sender()

	log.Printf("[USER:%d] Command: /start (@%s)", user.ID, user.Username)

	// Register or update user in database
	userData := &models.User{
		TelegramID: user.ID,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
	}

	if err := h.client.RegisterUser(userData); err != nil {
		log.Printf("[ERROR] Failed to register user %d: %v", user.ID, err)
		// Continue anyway - don't block user experience
	}

	// Check for referral parameter
	// Format: /start ref_123456789
	args := c.Args()
	if len(args) > 0 && len(args[0]) > 4 && args[0][:4] == "ref_" {
		referrerID := int64(0)
		fmt.Sscanf(args[0], "ref_%d", &referrerID)

		if referrerID > 0 && referrerID != user.ID {
			log.Printf("[REFERRAL] User %d referred by %d", user.ID, referrerID)
			resp, err := h.client.ProcessReferral(referrerID, user.ID)
			if err != nil {
				log.Printf("[ERROR] Failed to process referral: %v", err)
			} else if resp.Success {
				log.Printf("[REFERRAL] Success: %s", resp.Message)
			}
		}
	}

	// Send welcome message
	message := GetWelcomeMessage(user.FirstName)
	menu := GetMainMenu()

	return c.Send(message, menu)
}

// HandleTestDaily triggers a test daily reminder
func (h *Handler) HandleTestDaily(c tele.Context) error {
	user := c.Sender()
	// Mock streak for testing
	streak := 0
	message := GetDailyReminderMessage(user.FirstName, streak)
	menu := GetMainMenu()
	return c.Send(message, menu)
}

// HandleTestStreak triggers a test streak reminder
func (h *Handler) HandleTestStreak(c tele.Context) error {
	user := c.Sender()
	// Mock streak for testing
	streak := 5
	message := GetDailyReminderMessage(user.FirstName, streak)
	menu := GetMainMenu()
	return c.Send(message, menu)
}

// HandleDebugSchedule triggers the server to generate notification jobs
func (h *Handler) HandleDebugSchedule(c tele.Context) error {
	log.Println("DEBUG: Entering HandleDebugSchedule")
	if err := h.client.ScheduleNotifications(); err != nil {
		log.Printf("DEBUG: ScheduleNotifications failed: %v", err)
		return c.Send(fmt.Sprintf("❌ Failed to schedule: %v", err))
	}
	log.Println("DEBUG: ScheduleNotifications success")
	return c.Send("✅ Server triggered to schedule daily notifications!")
}

// GetMainMenu returns the main inline keyboard
func GetMainMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	btnPlay := menu.WebApp("🎮 Play DEC0D3 Game 🎮", &tele.WebApp{
		URL: "https://ushpuras.dev/DEC0D3/",
	})

	menu.Inline(
		menu.Row(btnPlay),
	)

	return menu
}
