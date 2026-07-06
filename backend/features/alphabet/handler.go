package alphabet

import (
	"net/http"
	"strconv"

	"spoony/database"

	"github.com/gin-gonic/gin"
)

// GetProgress returns the level IDs the current user has completed in
// "First Steps". The frontend already knows the level structure (titles,
// unlock order); this just reports which of the 10 are done.
func GetProgress(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := database.DB.Query(
		`SELECT level_id FROM alphabet_progress WHERE user_id = $1 AND completed = true`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	levelIDs := []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		levelIDs = append(levelIDs, id)
	}

	c.JSON(http.StatusOK, levelIDs)
}

// CompleteLevel records that the current user finished one First Steps level.
func CompleteLevel(c *gin.Context) {
	userID := c.GetInt("user_id")
	levelID, err := strconv.Atoi(c.Param("levelId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO alphabet_progress (user_id, level_id, completed, completed_at)
		VALUES ($1, $2, true, NOW())
		ON CONFLICT (user_id, level_id) DO UPDATE SET completed = true, completed_at = NOW()`,
		userID, levelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
