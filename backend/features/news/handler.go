package news

import (
	"net/http"

	"tma-boilerplate/database"

	"github.com/gin-gonic/gin"
)

// GetNews returns the Home feed, newest first.
func GetNews(c *gin.Context) {
	rows, err := database.DB.Query(
		`SELECT id, author, avatar, title, body, image, category, likes, created_at
		 FROM news_posts ORDER BY created_at DESC, id DESC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	posts := []Post{}
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Author, &p.Avatar, &p.Title, &p.Body,
			&p.Image, &p.Category, &p.Likes, &p.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		posts = append(posts, p)
	}

	c.JSON(http.StatusOK, posts)
}
