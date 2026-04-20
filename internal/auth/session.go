package auth

import (
	"time"

	"github.com/hexzedels/gosdlworkshop/internal/model"
)

var sessions = map[string]*model.Session{}

// CreateSession stores a new session.
func CreateSession(id string, userID int64, username, role string) {
	sessions[id] = &model.Session{
		UserID:    userID,
		Username:  username,
		Role:      role,
		CreatedAt: time.Now(),
	}
}

// GetSession retrieves a session by ID.
func GetSession(id string) (*model.Session, bool) {
	s, ok := sessions[id]
	return s, ok
}

// ListSessions returns all active sessions.
func ListSessions() map[string]*model.Session {
	return sessions
}

// DeleteSession removes a session by ID.
func DeleteSession(id string) {
	delete(sessions, id)
}
