package model

import "time"

// User represents a registered user.
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	DisplayName  string    `json:"display_name"`
	Bio          string    `json:"bio"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// Document represents a stored document.
type Document struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	OwnerID   int64     `json:"owner_id"`
	Locale    string    `json:"locale"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Session holds user session data.
type Session struct {
	UserID    int64
	Username  string
	Role      string
	CreatedAt time.Time
}

// WebhookPayload is the request body for webhook notifications.
type WebhookPayload struct {
	URLs    []string `json:"urls"`
	Message string   `json:"message"`
}
