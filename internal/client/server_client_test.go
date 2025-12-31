package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"decodeBot/internal/models"
)

func TestRegisterUser(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL
		if r.URL.Path != "/api/bot/register" {
			t.Errorf("Expected path /api/bot/register, got %s", r.URL.Path)
		}

		// Verify Method
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify Headers
		if r.Header.Get("X-Bot-Secret") != "test-secret" {
			t.Errorf("Expected X-Bot-Secret header to be test-secret")
		}

		// Verify Body
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			t.Errorf("Failed to decode body: %v", err)
		}
		if user.TelegramID != 123456 {
			t.Errorf("Expected TelegramID 123456, got %d", user.TelegramID)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	// Init Client
	client := NewServerClient(server.URL, "test-secret")

	// Execute
	user := &models.User{
		TelegramID: 123456,
		Username:   "test_bot_user",
		FirstName:  "Test",
	}
	err := client.RegisterUser(user)

	// Assert
	if err != nil {
		t.Errorf("RegisterUser returned error: %v", err)
	}
}

func TestProcessReferral(t *testing.T) {
	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL
		if r.URL.Path != "/api/bot/referral" {
			t.Errorf("Expected path /api/bot/referral, got %s", r.URL.Path)
		}

		// Verify Body
		var req models.ReferralRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode body: %v", err)
		}
		if req.ReferrerID != 111 || req.ReferredID != 222 {
			t.Errorf("Unexpected IDs in request: %v", req)
		}

		// Respond
		resp := models.ReferralResponse{
			Success:       true,
			ShardsAwarded: 20,
			Message:       "Referral successful",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Init Client
	client := NewServerClient(server.URL, "test-secret")

	// Execute
	resp, err := client.ProcessReferral(111, 222)

	// Assert
	if err != nil {
		t.Errorf("ProcessReferral returned error: %v", err)
	}
	if resp == nil || !resp.Success {
		t.Errorf("Expected successful response, got %v", resp)
	}
	if resp.ShardsAwarded != 20 {
		t.Errorf("Expected 20 shards, got %d", resp.ShardsAwarded)
	}
}
