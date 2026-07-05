package trainer

import (
	"net/http"
	"strconv"

	"spoony/database"

	"github.com/gin-gonic/gin"
)

// GetLevels returns all cities on the England route, enriched with the current
// user's progress. A level is unlocked if it's the first one or the previous
// level has been completed.
func GetLevels(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := database.DB.Query(`
		SELECT l.id, l.city, l.title_ru, l.description, l.emoji, l.order_index, l.pos_x, l.pos_y, l.game_type,
			COALESCE(p.completed, false), COALESCE(p.stars, 0)
		FROM levels l
		LEFT JOIN user_progress p ON p.level_id = l.id AND p.user_id = $1
		ORDER BY l.order_index ASC, l.id ASC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	levels := []Level{}
	prevCompleted := true // first level is always unlocked
	for rows.Next() {
		var l Level
		if err := rows.Scan(&l.ID, &l.City, &l.TitleRu, &l.Description, &l.Emoji,
			&l.OrderIndex, &l.PosX, &l.PosY, &l.GameType, &l.Completed, &l.Stars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		l.Unlocked = prevCompleted
		prevCompleted = l.Completed
		levels = append(levels, l)
	}

	c.JSON(http.StatusOK, levels)
}

// GetLevelCards returns the vocabulary items for a level. The frontend turns
// each item into a picture card + a word card for the "Find the pair" game.
func GetLevelCards(c *gin.Context) {
	levelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, level_id, word_en, word_ru, emoji FROM vocab_items
		 WHERE level_id = $1 ORDER BY id ASC`, levelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	items := []VocabItem{}
	for rows.Next() {
		var v VocabItem
		if err := rows.Scan(&v.ID, &v.LevelID, &v.WordEn, &v.WordRu, &v.Emoji); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		items = append(items, v)
	}

	if len(items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "level not found"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetLevelTheory returns the grammar cards for a "theory" level.
func GetLevelTheory(c *gin.Context) {
	levelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, level_id, order_index, title_ru, body_ru, example_en, example_ru
		 FROM theory_slides WHERE level_id = $1 ORDER BY order_index ASC, id ASC`, levelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	slides := []TheorySlide{}
	for rows.Next() {
		var s TheorySlide
		if err := rows.Scan(&s.ID, &s.LevelID, &s.OrderIndex, &s.TitleRu, &s.BodyRu,
			&s.ExampleEn, &s.ExampleRu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		slides = append(slides, s)
	}

	if len(slides) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no theory for this level"})
		return
	}

	c.JSON(http.StatusOK, slides)
}

// CompleteLevel records that the user finished a level, keeping the best star score.
func CompleteLevel(c *gin.Context) {
	userID := c.GetInt("user_id")
	levelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}

	var req CompleteLevelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Stars = 1
	}
	if req.Stars < 1 {
		req.Stars = 1
	}
	if req.Stars > 3 {
		req.Stars = 3
	}

	_, err = database.DB.Exec(`
		INSERT INTO user_progress (user_id, level_id, stars, completed, completed_at)
		VALUES ($1, $2, $3, true, NOW())
		ON CONFLICT (user_id, level_id) DO UPDATE SET
			stars = GREATEST(user_progress.stars, EXCLUDED.stars),
			completed = true,
			completed_at = NOW()`,
		userID, levelID, req.Stars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "level completed", "stars": req.Stars})
}

// GetUserStats returns simple aggregate stats for the Profile page.
func GetUserStats(c *gin.Context) {
	userID := c.GetInt("user_id")

	var totalLevels, completedLevels, totalStars, learnedWords int

	database.DB.QueryRow(`SELECT COUNT(*) FROM levels`).Scan(&totalLevels)
	database.DB.QueryRow(
		`SELECT COUNT(*), COALESCE(SUM(stars),0) FROM user_progress WHERE user_id = $1 AND completed = true`,
		userID).Scan(&completedLevels, &totalStars)
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM vocab_items v
		JOIN user_progress p ON p.level_id = v.level_id
		WHERE p.user_id = $1 AND p.completed = true`, userID).Scan(&learnedWords)

	c.JSON(http.StatusOK, gin.H{
		"total_levels":     totalLevels,
		"completed_levels": completedLevels,
		"total_stars":      totalStars,
		"learned_words":    learnedWords,
	})
}
