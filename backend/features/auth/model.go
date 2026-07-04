package auth

import "time"

type User struct {
	ID         int       `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email,omitempty"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
	Language   string    `json:"language"`
	Timezone   string    `json:"timezone,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Language  string `json:"language"`
	Timezone  string `json:"timezone"`
}
