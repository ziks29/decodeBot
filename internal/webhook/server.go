package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"decodeBot/internal/bot"

	tele "gopkg.in/telebot.v4"
)

// Server represents the webhook HTTP server
type Server struct {
	bot       *tele.Bot
	botSecret string
	port      string
}

// NewServer creates a new webhook server
func NewServer(bot *tele.Bot, port string) *Server {
	botSecret := os.Getenv("BOT_SECRET")
	return &Server{
		bot:       bot,
		botSecret: botSecret,
		port:      port,
	}
}

// NewUserRequest represents the request payload for new user notifications
type NewUserRequest struct {
	TelegramID int64  `json:"telegram_id"`
	FirstName  string `json:"first_name"`
}

// ReferralNotificationRequest represents the request payload for referral notifications
type ReferralNotificationRequest struct {
	ReferrerID   int64  `json:"referrer_id"`
	ReferredName string `json:"referred_name"`
}

// authenticateRequest checks if the request has the correct bot secret
func (s *Server) authenticateRequest(r *http.Request) bool {
	if s.botSecret == "" {
		// No secret configured, skip authentication
		return true
	}

	providedSecret := r.Header.Get("X-Bot-Secret")
	return providedSecret == s.botSecret
}

// handleNewUser handles incoming new user notifications from the backend
func (s *Server) handleNewUser(w http.ResponseWriter, r *http.Request) {
	// Authenticate request
	if !s.authenticateRequest(r) {
		log.Printf("[WEBHOOK] Unauthorized new user notification from %s", r.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req NewUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[WEBHOOK] Failed to parse new user request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.TelegramID == 0 {
		log.Printf("[WEBHOOK] Missing telegram_id in request")
		http.Error(w, "telegram_id is required", http.StatusBadRequest)
		return
	}

	log.Printf("[WEBHOOK] Received new user notification: TG ID %d (@%s)", req.TelegramID, req.FirstName)

	// Send welcome message
	message := bot.GetWelcomeMessage(req.FirstName)
	menu := bot.GetMainMenu()

	recipient := &tele.User{ID: req.TelegramID}
	if _, err := s.bot.Send(recipient, message, menu); err != nil {
		log.Printf("[WEBHOOK] Failed to send welcome message to user %d: %v", req.TelegramID, err)
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("[WEBHOOK] Successfully sent welcome message to user %d", req.TelegramID)

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Welcome message sent",
	})
	})
}

// handleReferral handles incoming referral notifications from the backend
func (s *Server) handleReferral(w http.ResponseWriter, r *http.Request) {
	// Authenticate request
	if !s.authenticateRequest(r) {
		log.Printf("[WEBHOOK] Unauthorized referral notification from %s", r.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req ReferralNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[WEBHOOK] Failed to parse referral request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.ReferrerID == 0 {
		log.Printf("[WEBHOOK] Missing referrer_id in request")
		http.Error(w, "referrer_id is required", http.StatusBadRequest)
		return
	}

	log.Printf("[WEBHOOK] Received referral notification: Referrer %d, Referred %s", req.ReferrerID, req.ReferredName)

	// Send notification message
	message := fmt.Sprintf("🚀 User **%s** just joined via your invite link!\n\n💎 You received +20 Shards!", req.ReferredName)
	recipient := &tele.User{ID: req.ReferrerID}
	
	if _, err := s.bot.Send(recipient, message, tele.ModeMarkdown); err != nil {
		log.Printf("[WEBHOOK] Failed to send referral message to user %d: %v", req.ReferrerID, err)
		// We perform a best-effort, so we don't return error to the server if the user blocked the bot
		// But we should log it.
	} else {
		log.Printf("[WEBHOOK] Successfully sent referral message to user %d", req.ReferrerID)
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Referral notification sent",
	})
}

// Start starts the webhook HTTP server
func (s *Server) Start() {
	http.HandleFunc("/webhook/new-user", s.handleNewUser)
	http.HandleFunc("/webhook/referral", s.handleReferral)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := ":" + s.port
	log.Printf("🌐 Webhook server starting on %s", addr)

	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Printf("❌ Webhook server failed: %v", err)
		}
	}()
}
