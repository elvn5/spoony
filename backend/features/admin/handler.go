package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"spoony/database"
	"spoony/features/news"

	"github.com/gin-gonic/gin"
)

// ---- Stats ----

func AdminGetStats(c *gin.Context) {
	var stats struct {
		TotalUsers      int `json:"total_users"`
		NewUsersToday   int `json:"new_users_today"`
		TelegramUsers   int `json:"telegram_users"`
		GuestUsers      int `json:"guest_users"`
		CompletedLevels int `json:"completed_levels"`
		TotalStars      int `json:"total_stars"`
		TotalNewsPosts  int `json:"total_news_posts"`
	}

	database.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers)
	database.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE created_at >= NOW() - INTERVAL '1 day'`).Scan(&stats.NewUsersToday)
	database.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE telegram_id IS NOT NULL`).Scan(&stats.TelegramUsers)
	database.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE guest_id IS NOT NULL`).Scan(&stats.GuestUsers)
	database.DB.QueryRow(`SELECT COUNT(*) FROM user_progress WHERE completed = true`).Scan(&stats.CompletedLevels)
	database.DB.QueryRow(`SELECT COALESCE(SUM(stars), 0) FROM user_progress WHERE completed = true`).Scan(&stats.TotalStars)
	database.DB.QueryRow(`SELECT COUNT(*) FROM news_posts`).Scan(&stats.TotalNewsPosts)

	c.JSON(http.StatusOK, stats)
}

// ---- Users ----

type userRow struct {
	ID         int    `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Username   string `json:"username"`
	Language   string `json:"language"`
	GuestID    string `json:"guest_id"`
	CreatedAt  string `json:"created_at"`
}

func AdminListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit
	search := "%" + c.Query("search") + "%"

	typeFilter := ""
	switch c.Query("type") {
	case "guest":
		typeFilter = "AND u.guest_id IS NOT NULL"
	case "telegram":
		typeFilter = "AND u.telegram_id IS NOT NULL"
	}

	rows, err := database.DB.Query(fmt.Sprintf(`
		SELECT u.id, COALESCE(u.telegram_id,0), u.first_name, u.last_name, u.username,
		       u.language, COALESCE(u.guest_id,''), u.created_at
		FROM users u
		WHERE (u.first_name ILIKE $1 OR u.username ILIKE $1 OR CAST(u.telegram_id AS TEXT) ILIKE $1) %s
		ORDER BY u.created_at DESC
		LIMIT $2 OFFSET $3`, typeFilter), search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []userRow
	for rows.Next() {
		var u userRow
		rows.Scan(&u.ID, &u.TelegramID, &u.FirstName, &u.LastName, &u.Username,
			&u.Language, &u.GuestID, &u.CreatedAt)
		users = append(users, u)
	}

	var total int
	database.DB.QueryRow(fmt.Sprintf(`
		SELECT COUNT(*) FROM users u
		WHERE (u.first_name ILIKE $1 OR u.username ILIKE $1 OR CAST(u.telegram_id AS TEXT) ILIKE $1) %s`, typeFilter),
		search).Scan(&total)

	c.JSON(http.StatusOK, gin.H{"users": users, "total": total, "page": page})
}

// AdminGetUser returns one user plus their progress across every level, so
// the admin can see (and then edit) exactly what a kid has and hasn't done.
func AdminGetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var user userRow
	err = database.DB.QueryRow(`
		SELECT id, COALESCE(telegram_id,0), first_name, last_name, username,
		       language, COALESCE(guest_id,''), created_at
		FROM users WHERE id = $1`, userID).Scan(
		&user.ID, &user.TelegramID, &user.FirstName, &user.LastName,
		&user.Username, &user.Language, &user.GuestID, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT l.id, l.city, l.title_ru, l.game_type, l.order_index,
		       COALESCE(p.completed, false), COALESCE(p.stars, 0)
		FROM levels l
		LEFT JOIN user_progress p ON p.level_id = l.id AND p.user_id = $1
		ORDER BY l.order_index ASC, l.id ASC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	type levelProgress struct {
		ID         int    `json:"id"`
		City       string `json:"city"`
		TitleRu    string `json:"title_ru"`
		GameType   string `json:"game_type"`
		OrderIndex int    `json:"order_index"`
		Completed  bool   `json:"completed"`
		Stars      int    `json:"stars"`
	}
	progress := []levelProgress{}
	for rows.Next() {
		var lp levelProgress
		if err := rows.Scan(&lp.ID, &lp.City, &lp.TitleRu, &lp.GameType, &lp.OrderIndex,
			&lp.Completed, &lp.Stars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		progress = append(progress, lp)
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "progress": progress})
}

// AdminUpdateUserProgress lets an admin correct a user's stars/completion on
// a single level (e.g. to unblock a kid stuck on a buggy level).
func AdminUpdateUserProgress(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	levelID, err := strconv.Atoi(c.Param("levelId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}

	var req UpdateProgressInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if req.Stars < 0 {
		req.Stars = 0
	}
	if req.Stars > 3 {
		req.Stars = 3
	}

	_, err = database.DB.Exec(`
		INSERT INTO user_progress (user_id, level_id, stars, completed, completed_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id, level_id) DO UPDATE SET
			stars = EXCLUDED.stars,
			completed = EXCLUDED.completed,
			completed_at = NOW()`,
		userID, levelID, req.Stars, req.Completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// AdminResetUserProgress deletes a user's progress on one level, re-locking
// it (and everything after it, since levels unlock sequentially).
func AdminResetUserProgress(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	levelID, err := strconv.Atoi(c.Param("levelId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid level id"})
		return
	}

	database.DB.Exec(`DELETE FROM user_progress WHERE user_id = $1 AND level_id = $2`, userID, levelID)
	c.JSON(http.StatusOK, gin.H{"ok": true})
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

// ---- Content (Home feed) ----

func AdminListNews(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, author, avatar, title, body, image, category, likes, created_at
		FROM news_posts ORDER BY created_at DESC, id DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	posts := []news.Post{}
	for rows.Next() {
		var p news.Post
		if err := rows.Scan(&p.ID, &p.Author, &p.Avatar, &p.Title, &p.Body,
			&p.Image, &p.Category, &p.Likes, &p.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		posts = append(posts, p)
	}
	c.JSON(http.StatusOK, posts)
}

func AdminCreateNews(c *gin.Context) {
	var req NewsPostInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author, title and body are required"})
		return
	}
	if req.Avatar == "" {
		req.Avatar = "🥄"
	}
	if req.Category == "" {
		req.Category = "news"
	}

	var id int
	err := database.DB.QueryRow(`
		INSERT INTO news_posts (author, avatar, title, body, image, category, likes)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		req.Author, req.Avatar, req.Title, req.Body, req.Image, req.Category, req.Likes,
	).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func AdminUpdateNews(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req NewsPostInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author, title and body are required"})
		return
	}
	if req.Avatar == "" {
		req.Avatar = "🥄"
	}
	if req.Category == "" {
		req.Category = "news"
	}

	res, err := database.DB.Exec(`
		UPDATE news_posts SET author=$1, avatar=$2, title=$3, body=$4, image=$5, category=$6, likes=$7
		WHERE id=$8`,
		req.Author, req.Avatar, req.Title, req.Body, req.Image, req.Category, req.Likes, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func AdminDeleteNews(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	database.DB.Exec(`DELETE FROM news_posts WHERE id=$1`, id)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
