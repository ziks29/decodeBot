package models

type User struct {
	TelegramID    int64  `json:"telegram_id"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	AllStreak     int    `json:"all_streak"`
	HexStreak     int    `json:"hex_streak"`
	WordStreak    int    `json:"word_streak"`
	NumericStreak int    `json:"numeric_streak"`
}

type UserProfile struct {
	TelegramID    int64  `json:"telegram_id"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	TotalGamesWon int    `json:"total_games_won"`
	CurrentStreak int    `json:"current_streak"`
	ShardBalance  int    `json:"shard_balance"`
	ReferralCount int    `json:"referral_count"`
	DailyStreak   int    `json:"daily_streak"`
	LastPlayedAt  string `json:"last_played_at"`
}

type ReferralRequest struct {
	ReferrerID int64 `json:"referrer_id"`
	ReferredID int64 `json:"referred_id"`
}

type ReferralResponse struct {
	Success       bool   `json:"success"`
	ShardsAwarded int    `json:"shards_awarded"`
	Message       string `json:"message"`
}
