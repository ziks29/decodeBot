# DEC0D3 Telegram Bot

A Telegram bot companion for the DEC0D3 cipher puzzle game. The bot provides an entry point for users to discover and launch the Mini App.

## 🎯 Purpose

- Welcome new users to DEC0D3
- Provide easy access to the Mini App
- Display user statistics and game information
- Send daily challenge reminders (optional)

## 🛠️ Tech Stack

- **Language:** Go 1.24+
- **Bot Library:** [telebot](https://gopkg.in/telebot.v4) v4 (beta.7)
- **Scheduler:** robfig/cron v3 for daily notifications
- **HTTP Client:** Native `net/http`
- **Environment:** `godotenv`
- **Container:** Docker


## 📦 Installation

### Prerequisites

- Go 1.24 or higher
- Telegram Bot Token (get from [@BotFather](https://t.me/botfather))
- Running instance of `decodeServer`

### Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Configure environment:**
   Create a `.env` file:
   ```env
   BOT_TOKEN=your_telegram_bot_token_here
   BOT_USERNAME=your_bot_username
   SERVER_URL=http://localhost:8081
   MINI_APP_URL=https://ushpuras.dev/DEC0D3/
   DEBUG=true
   ```

3. **Run locally:**
   ```bash
   go run cmd/bot/main.go
   ```

## 🐳 Docker

### Build Image

```bash
docker build -t decode-bot .
```

### Run Container

```bash
docker run -d --name decode-bot --env-file .env decode-bot
```

### Docker Compose

```bash
# Run both bot and server
docker-compose up -d
```

## 📝 Bot Commands

| Command | Description | Status |
|---------|-------------|--------|
| `/start` | Welcome message + launch Mini App + referral tracking | ✅ Implemented |
| `/help` | Game instructions and features | 🚧 Coming Soon |
| `/stats` | Personal game statistics | 🚧 Coming Soon |
| `/daily` | Today's daily challenge info | 🚧 Coming Soon |

## 🔔 Automated Features

- **Daily Reminders** - 9:00 AM reminder for daily challenges
- **Streak Reminders** - 8:00 PM reminder for users with active streaks
- **Referral System** - +20 shards for both referrer and referred user


## 🔧 Development

### Project Structure

```
decodeBot/
├── cmd/
│   └── bot/
│       └── main.go              # Entry point
├── internal/
│   ├── bot/
│   │   ├── handlers.go          # Command handlers
│   │   ├── keyboards.go         # Inline keyboards
│   │   └── messages.go          # Message templates
│   ├── client/
│   │   └── server_client.go     # API client
│   ├── config/
│   │   └── config.go            # Configuration
│   └── models/
│       └── user.go              # Data models
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o bot cmd/bot/main.go
```

## 🔐 Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `BOT_TOKEN` | Telegram bot token from BotFather | ✅ | - |
| `BOT_USERNAME` | Bot username (without @) | ✅ | - |
| `SERVER_URL` | Base URL of decodeServer | ✅ | `http://localhost:8081` |
| `MINI_APP_URL` | URL of the Mini App | ✅ | - |
| `DEBUG` | Enable debug logging | ❌ | `false` |
| `LOG_LEVEL` | Log level (info, debug, error) | ❌ | `info` |

## 🌐 Server API Integration

The bot communicates with `decodeServer` using these endpoints:

- `GET /api/health` - Check server health
- `POST /api/bot/register` - Register new user
- `GET /api/bot/stats/:telegramId` - Get user statistics
- `GET /api/leaderboard` - Fetch leaderboard
- `GET /api/daily-challenge` - Get today's challenge

## 📊 Monitoring

The bot includes structured logging:

```
[USER:123456789] Command: /start
[API] GET /api/health -> 200 OK
[ERROR] Failed to send message: connection timeout
```

## 🚀 Deployment

### Production Checklist

- [ ] Set production `BOT_TOKEN`
- [ ] Configure production `SERVER_URL`
- [ ] Set `DEBUG=false`
- [ ] Enable HTTPS for webhooks (optional)
- [ ] Set up log aggregation
- [ ] Configure restart policies
- [ ] Set up health checks

### Webhook Mode (Optional)

For production, you can use webhooks instead of polling:

```bash
# Set webhook
curl -X POST https://api.telegram.org/bot<TOKEN>/setWebhook \
  -d url=https://yourdomain.com/webhook
```

## 📚 Resources

- [Telegram Bot API Documentation](https://core.telegram.org/bots/api)
- [go-telegram-bot-api Documentation](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5)
- [Telegram Mini Apps Guide](https://core.telegram.org/bots/webapps)
- [DEC0D3 Main Repository](../README.md)

## 🐛 Troubleshooting

### Bot Not Responding

1. Check bot token is correct
2. Verify server is running
3. Check logs for errors
4. Test with `/start` command

### Mini App Not Opening

1. Verify `MINI_APP_URL` is correct
2. Check bot username matches
3. Ensure Mini App is deployed
4. Test URL in browser first

### API Errors

1. Check `decodeServer` is running
2. Verify `SERVER_URL` is accessible
3. Check network connectivity
4. Review server logs

## 📄 License

All rights reserved. This is proprietary software for the DEC0D3 project.

---

**Status:** 🚧 Under Development

For the full implementation plan, see [TELEGRAM_BOT_PLAN.md](../TELEGRAM_BOT_PLAN.md)
