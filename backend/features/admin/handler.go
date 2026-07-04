package admin

import (
	"net/http"
	"strconv"

	"tma-boilerplate/database"

	"github.com/gin-gonic/gin"
)

// ---- Stats ----

func AdminGetStats(c *gin.Context) {
	var stats struct {
		TotalUsers    int `json:"total_users"`
		NewUsersToday int `json:"new_users_today"`
	}

	database.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers)
	database.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE created_at >= NOW() - INTERVAL '1 day'`).Scan(&stats.NewUsersToday)

	c.JSON(http.StatusOK, stats)
}

// ---- Users ----

func AdminListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 20
	offset := (page - 1) * limit
	search := "%" + c.Query("search") + "%"

	rows, err := database.DB.Query(`
		SELECT u.id, u.telegram_id, u.first_name, u.last_name, u.username,
		       u.language, u.created_at
		FROM users u
		WHERE u.first_name ILIKE $1 OR u.username ILIKE $1 OR CAST(u.telegram_id AS TEXT) ILIKE $1
		ORDER BY u.created_at DESC
		LIMIT $2 OFFSET $3`, search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type UserRow struct {
		ID         int    `json:"id"`
		TelegramID int64  `json:"telegram_id"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Username   string `json:"username"`
		Language   string `json:"language"`
		CreatedAt  string `json:"created_at"`
	}

	var users []UserRow
	for rows.Next() {
		var u UserRow
		rows.Scan(&u.ID, &u.TelegramID, &u.FirstName, &u.LastName, &u.Username,
			&u.Language, &u.CreatedAt)
		users = append(users, u)
	}

	var total int
	database.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE first_name ILIKE $1 OR username ILIKE $1 OR CAST(telegram_id AS TEXT) ILIKE $1`, search).Scan(&total)

	c.JSON(http.StatusOK, gin.H{"users": users, "total": total, "page": page})
}

func AdminDeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	database.DB.Exec(`DELETE FROM users WHERE id=$1`, userID)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
