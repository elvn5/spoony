package news

import "time"

// Post is a single item in the Home feed (Facebook-style timeline).
type Post struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`
	Avatar    string    `json:"avatar"` // emoji used as avatar
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Image     string    `json:"image"` // big emoji / illustration
	Category  string    `json:"category"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
}
