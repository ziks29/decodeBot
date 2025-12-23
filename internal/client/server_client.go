package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"decodeBot/internal/models"
)

type ServerClient struct {
	baseURL    string
	botSecret  string
	httpClient *http.Client
	maxRetries int
	retryDelay time.Duration
}

func NewServerClient(baseURL, botSecret string) *ServerClient {
	client := &ServerClient{
		baseURL:   baseURL,
		botSecret: botSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		maxRetries: 3,
		retryDelay: 1 * time.Second,
	}
	log.Printf("[CLIENT] Initialized ServerClient with URL: %s", baseURL)
	return client
}

// doWithRetry executes an HTTP request with retry logic
func (c *ServerClient) doWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)

		// Success - return immediately
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// Last attempt - return error
		if attempt == c.maxRetries {
			if err != nil {
				log.Printf("[CLIENT] Request failed after %d attempts: %v", c.maxRetries+1, err)
				return nil, err
			}
			log.Printf("[CLIENT] Request failed with status %d after %d attempts", resp.StatusCode, c.maxRetries+1)
			return resp, nil
		}

		// Calculate delay with exponential backoff
		delay := c.retryDelay * time.Duration(1<<uint(attempt))
		if delay > 5*time.Second {
			delay = 5 * time.Second
		}

		if err != nil {
			log.Printf("[CLIENT] Request failed (attempt %d/%d): %v, retrying in %v...", attempt+1, c.maxRetries+1, err, delay)
		} else {
			log.Printf("[CLIENT] Request returned %d (attempt %d/%d), retrying in %v...", resp.StatusCode, attempt+1, c.maxRetries+1, delay)
			resp.Body.Close()
		}

		time.Sleep(delay)
	}

	return resp, err
}

// HealthCheck checks if the server is running
func (c *ServerClient) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/api/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server health check failed: %d", resp.StatusCode)
	}

	return nil
}

// RegisterUser registers a new user from the bot
func (c *ServerClient) RegisterUser(user *models.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/bot/register", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.botSecret != "" {
		req.Header.Set("X-Bot-Secret", c.botSecret)
	}

	resp, err := c.doWithRetry(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to register user: %d - %s", resp.StatusCode, string(body))
	}

	log.Printf("[API] User registered: %d (@%s)", user.TelegramID, user.Username)
	return nil
}

// GetUserProfile fetches user profile data
func (c *ServerClient) GetUserProfile(telegramID int64) (*models.UserProfile, error) {
	url := fmt.Sprintf("%s/api/bot/stats/%d", c.baseURL, telegramID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.botSecret != "" {
		req.Header.Set("X-Bot-Secret", c.botSecret)
	}

	resp, err := c.doWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user profile: %d - %s", resp.StatusCode, string(body))
	}

	var profile models.UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// ProcessReferral processes a referral and awards shards
func (c *ServerClient) ProcessReferral(referrerID, referredID int64) (*models.ReferralResponse, error) {
	req := models.ReferralRequest{
		ReferrerID: referrerID,
		ReferredID: referredID,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/api/bot/referral", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.botSecret != "" {
		httpReq.Header.Set("X-Bot-Secret", c.botSecret)
	}

	resp, err := c.doWithRetry(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var refResp models.ReferralResponse
	if err := json.NewDecoder(resp.Body).Decode(&refResp); err != nil {
		return nil, err
	}

	return &refResp, nil
}

// NotificationJob represents a scheduled notification
type NotificationJob struct {
	ID         uint  `json:"id"`
	UserID     uint  `json:"user_id"`
	TelegramID int64 `json:"telegram_id"` // Added this to response in server? Wait, server returns job which has UserID. I need to join with User table?
	// The server model NotificationJob has UserID. The pending notifications endpoint should probably return user info too?
	// Let's check server service GetPendingNotifications. It returns []models.NotificationJob.
	// DOES models.NotificationJob include User data?
	// In server models.go: NotificationJob has UserID but no User association defined in struct yet.
	// I need to update server model to include User association or manually fetch it.
	// Let's assume for now I will fix server to Preload User.
	// So Client struct needs User info.
	Type        string       `json:"type"`
	ScheduledAt time.Time    `json:"scheduled_at"`
	Status      string       `json:"status"`
	User        *models.User `json:"user"` // Nested user object
}

// ScheduleNotifications triggers manual scheduling on server
func (c *ServerClient) ScheduleNotifications() error {
	url := c.baseURL + "/api/bot/notifications/schedule"
	log.Printf("[CLIENT] POST %s", url)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	if c.botSecret != "" {
		req.Header.Set("X-Bot-Secret", c.botSecret)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[CLIENT] Request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("[CLIENT] Response status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to schedule notifications: %d", resp.StatusCode)
	}
	return nil
}

// GetPendingNotifications fetches pending jobs from server
func (c *ServerClient) GetPendingNotifications(limit int) ([]NotificationJob, error) {
	url := fmt.Sprintf("%s/api/bot/notifications/pending?limit=%d", c.baseURL, limit)
	log.Printf("[CLIENT] GET %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if c.botSecret != "" {
		req.Header.Set("X-Bot-Secret", c.botSecret)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[CLIENT] Request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("[CLIENT] Response status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get pending notifications: %d", resp.StatusCode)
	}

	var jobs []NotificationJob
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

// UpdateJobStatus updates the status of a job
func (c *ServerClient) UpdateJobStatus(jobID uint, status string) error {
	payload := map[string]string{"status": status}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/bot/notifications/%d", c.baseURL, jobID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.botSecret != "" {
		req.Header.Set("X-Bot-Secret", c.botSecret)
	}

	resp, err := c.doWithRetry(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update job status: %d", resp.StatusCode)
	}
	return nil
}

// GetUserStats fetches user statistics from the server
func (c *ServerClient) GetUserStats() (*models.UserStats, error) {
	url := c.baseURL + "/api/bot/stats"
	log.Printf("[CLIENT] GET %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if c.botSecret != "" {
		req.Header.Set("X-Bot-Secret", c.botSecret)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("[CLIENT] Request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("[CLIENT] Response status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user stats: %d", resp.StatusCode)
	}

	var stats models.UserStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
