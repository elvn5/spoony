package alphabet

import "time"

// Progress records that a user finished one of the 10 "First Steps"
// levels (4 base levels + 6 letter-combo groups, defined in frontend data).
type Progress struct {
	UserID      int       `json:"user_id"`
	LevelID     int       `json:"level_id"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `json:"completed_at"`
}
