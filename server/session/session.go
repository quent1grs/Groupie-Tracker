package salon

import (
	"groupietracker/server/user"
)

type Salon struct {
	SalonUsers []user.User
	SalonID    string
	// SalonGame  Game
	// SalonChat        Chat
	// SalonLeaderboard Leaderboard
}

func GetSession() {

}
