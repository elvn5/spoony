package trainer

import "time"

// Level is one city on the England route. Completing it unlocks the next one.
type Level struct {
	ID          int    `json:"id"`
	City        string `json:"city"`
	TitleRu     string `json:"title_ru"`
	Description string `json:"description"`
	Emoji       string `json:"emoji"`
	OrderIndex  int    `json:"order_index"`
	PosX        int    `json:"pos_x"` // 0..100 position on the map
	PosY        int    `json:"pos_y"`
	GameType    string `json:"game_type"` // "match" (find the pair) or "word_build" (assemble a word)

	// Per-user fields (filled in when fetched for an authenticated user).
	Completed bool `json:"completed"`
	Stars     int  `json:"stars"`
	Unlocked  bool `json:"unlocked"`
}

// VocabItem is a vocabulary entry. Each item produces a pair of cards in the
// "Find the pair" game: one card shows the emoji (picture), the other the word.
type VocabItem struct {
	ID      int    `json:"id"`
	LevelID int    `json:"level_id"`
	WordEn  string `json:"word_en"`
	WordRu  string `json:"word_ru"`
	Emoji   string `json:"emoji"`
}

// UserProgress records completion of a level by a user.
type UserProgress struct {
	UserID      int       `json:"user_id"`
	LevelID     int       `json:"level_id"`
	Stars       int       `json:"stars"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `json:"completed_at"`
}

// CompleteLevelRequest is the body sent when a user finishes a level.
type CompleteLevelRequest struct {
	Stars int `json:"stars"`
}
