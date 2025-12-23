# Backend Endpoints for Bot Integration

## Overview

These are the new endpoints needed in `decodeServer` to support the Telegram bot functionality.

## Endpoints to Add

### 1. POST /api/bot/register

**Purpose:** Register or update a user from bot interaction

**Request:**
```json
{
  "telegram_id": 123456789,
  "username": "user123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "is_new_user": true
}
```

**Implementation Notes:**
- Check if user exists by `telegram_id`
- If new user, create profile with default values
- If existing, update username/name if changed
- Return `is_new_user: true` for analytics

**Handler Location:** `decodeServer/internal/handlers/bot.go`

---

### 2. GET /api/bot/stats/:telegramId

**Purpose:** Get user profile and stats for bot display

**Response:**
```json
{
  "telegram_id": 123456789,
  "username": "user123",
  "first_name": "John",
  "total_games_won": 42,
  "current_streak": 5,
  "daily_streak": 3,
  "shard_balance": 150,
  "referral_count": 2,
  "last_played_at": "2025-12-22T10:30:00Z"
}
```

**Implementation Notes:**
- Similar to `/api/profile` but optimized for bot
- Include referral count
- Include streak information
- No authentication needed (bot has special access)

**Handler Location:** `decodeServer/internal/handlers/bot.go`

---

### 3. POST /api/bot/referral

**Purpose:** Process referral and award shards to both users

**Request:**
```json
{
  "referrer_id": 123456789,
  "referred_id": 987654321
}
```

**Response Success:**
```json
{
  "success": true,
  "shards_awarded": 20,
  "message": "Both users received +20 shards!"
}
```

**Response Error (Already Referred):**
```json
{
  "success": false,
  "shards_awarded": 0,
  "message": "This user has already been referred"
}
```

**Implementation Notes:**
- Validate that `referrer_id != referred_id` (no self-referrals)
- Check if `referred_id` has already been referred by anyone
- If valid:
  - Add +20 shards to referrer
  - Add +20 shards to referred user
  - Create referral record in database
  - Increment referrer's `referral_count`
- Use database transaction for atomicity

**Handler Location:** `decodeServer/internal/handlers/bot.go`

**Database Schema Needed:**
```sql
CREATE TABLE IF NOT EXISTS referrals (
    id SERIAL PRIMARY KEY,
    referrer_id BIGINT NOT NULL,
    referred_id BIGINT NOT NULL UNIQUE,  -- Each user can only be referred once
    shards_awarded INT DEFAULT 20,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (referrer_id) REFERENCES profiles(telegram_id),
    FOREIGN KEY (referred_id) REFERENCES profiles(telegram_id)
);

CREATE INDEX idx_referrals_referrer ON referrals(referrer_id);
```

---

### 4. GET /api/bot/active-users (Optional)

**Purpose:** Get list of users for daily notifications

**Query Parameters:**
- `limit` - Max users to return (default: 1000)
- `offset` - Pagination offset (default: 0)
- `min_games` - Min games played to be considered active (default: 1)

**Response:**
```json
{
  "users": [
    {
      "telegram_id": 123456789,
      "first_name": "John",
      "current_streak": 5,
      "last_played_at": "2025-12-22T10:30:00Z"
    }
  ],
  "total": 1250,
  "limit": 1000,
  "offset": 0
}
```

**Implementation Notes:**
- Return users who have played at least `min_games`
- Order by `last_played_at DESC`
- Use pagination for large user bases
- Only include active users (played in last 30 days)

**Handler Location:** `decodeServer/internal/handlers/bot.go`

---

## Middleware Considerations

### Bot Authentication

The bot needs special access without regular Telegram auth. Options:

**Option 1: Shared Secret**
```go
func BotAuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        botSecret := c.Get("X-Bot-Secret")
        expectedSecret := os.Getenv("BOT_SECRET")
        
        if botSecret == "" || botSecret != expectedSecret {
            return c.Status(401).JSON(fiber.Map{
                "error": "Unauthorized bot request",
            })
        }
        
        return c.Next()
    }
}
```

**Option 2: IP Whitelist**
```go
func BotAuthMiddleware() fiber.Handler {
    allowedIPs := []string{"127.0.0.1", "YOUR_BOT_SERVER_IP"}
    
    return func(c *fiber.Ctx) error {
        clientIP := c.IP()
        
        for _, ip := range allowedIPs {
            if clientIP == ip {
                return c.Next()
            }
        }
        
        return c.Status(403).JSON(fiber.Map{
            "error": "IP not whitelisted",
        })
    }
}
```

**Recommended:** Use Option 1 (Shared Secret) for flexibility

---

## Routes Setup

Add to `decodeServer/main.go`:

```go
// Bot routes (with bot authentication)
botAPI := api.Group("/bot")
botAPI.Use(middleware.BotAuth()) // Custom middleware

botAPI.Post("/register", handlers.BotRegisterUser)
botAPI.Get("/stats/:telegramId", handlers.BotGetUserStats)
botAPI.Post("/referral", handlers.BotProcessReferral)
botAPI.Get("/active-users", handlers.BotGetActiveUsers) // Optional
```

---

## Database Changes

### Add to profiles table:
```sql
ALTER TABLE profiles ADD COLUMN referral_count INT DEFAULT 0;
```

### Create referrals table:
```sql
CREATE TABLE IF NOT EXISTS referrals (
    id SERIAL PRIMARY KEY,
    referrer_id BIGINT NOT NULL,
    referred_id BIGINT NOT NULL UNIQUE,
    shards_awarded INT DEFAULT 20,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (referrer_id) REFERENCES profiles(telegram_id),
    FOREIGN KEY (referred_id) REFERENCES profiles(telegram_id)
);

CREATE INDEX idx_referrals_referrer ON referrals(referrer_id);
CREATE INDEX idx_referrals_referred ON referrals(referred_id);
```

---

## Environment Variables

Add to `decodeServer/.env`:

```env
# Bot Configuration
BOT_SECRET=your_secure_random_secret_here
BOT_ENABLED=true
```

---

## Testing the Endpoints

### Test User Registration
```bash
curl -X POST http://localhost:8081/api/bot/register \
  -H "Content-Type: application/json" \
  -H "X-Bot-Secret: your_secret" \
  -d '{
    "telegram_id": 123456789,
    "username": "testuser",
    "first_name": "Test",
    "last_name": "User"
  }'
```

### Test User Stats
```bash
curl http://localhost:8081/api/bot/stats/123456789 \
  -H "X-Bot-Secret: your_secret"
```

### Test Referral
```bash
curl -X POST http://localhost:8081/api/bot/referral \
  -H "Content-Type: application/json" \
  -H "X-Bot-Secret: your_secret" \
  -d '{
    "referrer_id": 123456789,
    "referred_id": 987654321
  }'
```

---

## Priority

1. **High Priority:** `/api/bot/register` and `/api/bot/referral`
   - Required for basic bot functionality
   
2. **Medium Priority:** `/api/bot/stats/:telegramId`
   - Needed for future commands like `/stats`
   
3. **Low Priority:** `/api/bot/active-users`
   - Only needed for Automated notifications
   - Can hardcode user list initially

---

## Next Steps

1. Create `decodeServer/internal/handlers/bot.go`
2. Create `decodeServer/internal/middleware/bot_auth.go`
3. Add database migration for referrals table
4. Update `main.go` with bot routes
5. Test locally with bot
6. Deploy and verify

---

**Estimated Time:** 2-3 hours for full implementation
