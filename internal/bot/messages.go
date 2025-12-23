package bot

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetWelcomeMessage returns the welcome message for /start command
func GetWelcomeMessage(firstName string) string {
	return fmt.Sprintf(`🔐 Welcome to DEC0D3, %s!

DEC0D3 is a cyber-themed cipher puzzle game where you decode secret patterns.

🎯 Game Variants:
• HEX - Decode 4-digit color codes
• NUMERIC - Guess 5-digit numbers
• WORD - Find 5-letter English words

✨ Features:
• 📅 Daily challenges with streak tracking
• 🏆 Global leaderboards
• 💎 Earn shards, get AI hints
• 🤖 Powered by Gemini AI
• 🎁 Invite friends and earn +20 shards per referral!

Ready to test your decoding skills?
Click the button below to start playing! 👇`, firstName)
}

// Cyberpunk messages for users WITH streaks
var streakMessages = []string{
	// Message 1: System Alert
	`⚡ SYSTEM BREACH DETECTED

Agent %s, your neural link has been active for %d cycles.

New encrypted data packets await extraction. Daily security protocols require immediate attention.

Continue your streak. Decrypt the codes. 🔐`,

	// Message 2: Network Status
	`🌐 NETWORK STATUS: ACTIVE

%s | Streak: %d days | Status: ELITE

The grid never sleeps. Today's transmission contains critical intel. Your pattern recognition skills are needed.

Access the mainframe now ⚡`,

	// Message 3: AI Companion
	`🤖 NEURAL AI REPORT

Hello %s. You've maintained cognitive sync for %d consecutive sessions.

Today's challenge matrix is loaded. The algorithms are waiting for your input. Don't let your streak flatline.

Engage protocols 🧠`,

	// Message 4: Urgent Transmission
	`📡 INCOMING: Priority Signal

%s, you're %d days deep in the simulation.

Today's ciphertext just dropped. The corporation doesn't rest, and neither should you. Decode before the window closes.

Stay connected 🔴`,

	// Message 5: Hacker Collective
	`👾 COLLECTIVE BROADCAST

%s - %d day operative streak recorded.

New targets identified. Your decryption skills put you in the top tier. The puzzles won't solve themselves, agent.

Jack in 🎮`,

	// Message 6: Memory Fragment
	`💾 MEMORY FRAGMENT DETECTED

Agent %s, %d continuous days logged in the archives.

Fresh data corruption needs your expertise. The hex, numeric, and word layers all require your touch. Time-sensitive.

Initialize sequence 🔍`,

	// Message 7: Glitch Aesthetic
	`█▀▀ █▀█ █▀▄ █▀▀   █▀▄ █▀█ █▀█ █▀█
█▄▄ █▄█ █▄▀ ██▄   █▄▀ █▀▄ █▄█ █▀▀

%s // STREAK: %d DAYS

New patterns emerged in the noise. Your presence is required for analysis. Don't break the chain.

>_ Execute now`,

	// Message 8: Surveillance Warning
	`👁️ SURVEILLANCE DETECTED

%s, you've been tracked for %d days straight.

They're watching your moves. Today's encrypted challenges are your only defense. Stay sharp, stay decoding, stay ahead.

Don't go dark now 🌙`,

	// Message 9: Crypto Mining
	`⛏️ CRYPTO MINING STATUS

Miner: %s | Uptime: %d days

Fresh hash puzzles ready for processing. Your neural network performance has been exceptional. Keep the computational power flowing.

Mine the codes 💎`,

	// Message 10: Reality Glitch
	`🔮 REALITY.EXE UNSTABLE

%s, the simulation recognizes your %d-day presence.

Today's glitches in the matrix reveal new patterns. Decode them before they vanish. The red pill is daily challenges.

Enter the void ⚡`,
}

// Cyberpunk messages for users WITHOUT streaks
var noStreakMessages = []string{
	// Message A: New Agent Onboarding
	`🌐 INITIALIZATION SEQUENCE

Welcome, Agent %s.

The network has registered your presence. Daily operations begin now. Your first mission: decrypt today's data streams.

Start your streak. Prove your worth 🔐`,

	// Message B: System Reboot
	`⚡ NEURAL LINK: RECONNECTING

%s, systems are back online.

You've been offline too long. The codes are piling up. Today's your chance to re-establish your streak and climb the ranks.

Reboot complete. Deploy now 🤖`,

	// Message C: Recruitment
	`📡 RECRUITMENT: ACTIVE

The collective needs decoders like you, %s.

Fresh intel just hit the network. HEX signatures, NUMERIC sequences, WORD ciphers—all waiting. Start your operation today.

Join the elite 👾`,

	// Message D: Challenge Issued
	`💾 NEW CHALLENGER DETECTED

%s, your skills haven't been forgotten.

The system remembers your last session. Today's challenges are calling. Build your streak from zero. Show them you're still sharp.

Accept protocol? Y/N_ 🔍`,

	// Message E: Data Leak
	`🔴 DATA LEAK IN PROGRESS

%s, unauthorized access detected in sector 7.

Only elite decoders can patch the breach. Today's puzzles hold the key. Start your streak and secure the network.

Time is running out ⚡`,
}

// GetDailyReminderMessage returns a random cyberpunk-themed daily reminder message
func GetDailyReminderMessage(firstName string, currentStreak int) string {
	if currentStreak > 0 {
		// Random message from streak messages
		idx := rand.Intn(len(streakMessages))
		return fmt.Sprintf(streakMessages[idx], firstName, currentStreak)
	}

	// Random message from no-streak messages
	idx := rand.Intn(len(noStreakMessages))
	return fmt.Sprintf(noStreakMessages[idx], firstName)
}

// GetStreakStatsMessage returns daily streak statistics
func GetStreakStatsMessage(profile interface{}) string {
	// We'll implement this when we have the profile structure from server
	return "📊 Your streak stats will appear here!"
}
