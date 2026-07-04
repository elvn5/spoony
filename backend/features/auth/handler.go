package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"tma-boilerplate/config"
	"tma-boilerplate/database"
	"tma-boilerplate/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TelegramLoginRequest struct {
	InitData string `json:"init_data" binding:"required"`
}

func TelegramLogin(c *gin.Context) {
	var req TelegramLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tgUser struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
		PhotoURL  string `json:"photo_url"`
	}

	vals, err := VerifyTelegramInitData(req.InitData)
	if err != nil {
		// Dev convenience: outside Telegram (e.g. browser at localhost) the
		// initData can't be verified. Sign in as a demo kid so the whole app
		// stays testable. Disabled in production.
		if config.App.Env != "production" {
			tgUser.ID = 777000777
			tgUser.FirstName = "Demo"
			tgUser.LastName = "Kid"
			tgUser.Username = "demo_kid"
			user, uerr := upsertUser(tgUser.ID, tgUser.Username, tgUser.FirstName, tgUser.LastName, tgUser.PhotoURL)
			if uerr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create demo user"})
				return
			}
			token, terr := generateJWT(user.ID)
			if terr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid telegram data: " + err.Error()})
		return
	}

	// Parse user from Telegram's JSON object in initData
	if userJSON := vals.Get("user"); userJSON != "" {
		if err := json.Unmarshal([]byte(userJSON), &tgUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user data"})
			return
		}
	} else {
		// Fallback: read individual fields (older initData format)
		tgUser.ID, _ = strconv.ParseInt(vals.Get("id"), 10, 64)
		tgUser.FirstName = vals.Get("first_name")
		tgUser.LastName = vals.Get("last_name")
		tgUser.Username = vals.Get("username")
	}

	if tgUser.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing telegram user id"})
		return
	}

	// Upsert user
	user, err := upsertUser(tgUser.ID, tgUser.Username, tgUser.FirstName, tgUser.LastName, tgUser.PhotoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create/update user"})
		return
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

type GuestLoginRequest struct {
	GuestID string `json:"guest_id"`
	Name    string `json:"name"`
}

// GuestLogin lets a user sign in without Telegram — as on a regular website.
// The browser sends a persistent guest_id (stored in localStorage) so a
// returning visitor keeps the same account and progress. Works in production.
func GuestLogin(c *gin.Context) {
	var req GuestLoginRequest
	_ = c.ShouldBindJSON(&req)

	req.GuestID = strings.TrimSpace(req.GuestID)
	if req.GuestID == "" {
		req.GuestID = "g_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = "Гость"
	}

	user, err := upsertGuestUser(req.GuestID, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create guest"})
		return
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user, "guest_id": req.GuestID})
}

func Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func GetMe(c *gin.Context) {
	userID := c.GetInt("user_id")

	var user User
	err := database.DB.QueryRow(`SELECT id, COALESCE(telegram_id,0), COALESCE(username,''), COALESCE(email,''), first_name, last_name,
		COALESCE(avatar_url,''), COALESCE(language,''), COALESCE(timezone,''), created_at, updated_at
		FROM users WHERE id = $1`, userID).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.Email,
		&user.FirstName, &user.LastName, &user.AvatarURL,
		&user.Language, &user.Timezone, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec(`UPDATE users SET first_name=$1, last_name=$2, language=$3, timezone=$4, updated_at=NOW()
		WHERE id=$5`, req.FirstName, req.LastName, req.Language, req.Timezone, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

func upsertUser(telegramID int64, username, firstName, lastName, avatarURL string) (*User, error) {
	var user User
	err := database.DB.QueryRow(`
		INSERT INTO users (telegram_id, username, first_name, last_name, avatar_url)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (telegram_id) DO UPDATE SET
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW()
		RETURNING id, telegram_id, username, COALESCE(email,''), first_name, last_name,
			COALESCE(avatar_url,''), COALESCE(language,'en'), COALESCE(timezone,''), created_at, updated_at`,
		telegramID, username, firstName, lastName, avatarURL,
	).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.Email,
		&user.FirstName, &user.LastName, &user.AvatarURL,
		&user.Language, &user.Timezone, &user.CreatedAt, &user.UpdatedAt,
	)
	return &user, err
}

func upsertGuestUser(guestID, firstName string) (*User, error) {
	username := "guest_" + guestID
	if len(username) > 32 {
		username = username[:32]
	}
	var user User
	err := database.DB.QueryRow(`
		INSERT INTO users (guest_id, username, first_name, last_name, language)
		VALUES ($1, $2, $3, '', 'ru')
		ON CONFLICT (guest_id) DO UPDATE SET
			first_name = EXCLUDED.first_name,
			updated_at = NOW()
		RETURNING id, COALESCE(telegram_id,0), COALESCE(username,''), COALESCE(email,''), first_name, last_name,
			COALESCE(avatar_url,''), COALESCE(language,'ru'), COALESCE(timezone,''), created_at, updated_at`,
		guestID, username, firstName,
	).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.Email,
		&user.FirstName, &user.LastName, &user.AvatarURL,
		&user.Language, &user.Timezone, &user.CreatedAt, &user.UpdatedAt,
	)
	return &user, err
}

func generateJWT(userID int) (string, error) {
	claims := middleware.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.App.JWTExpirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.App.JWTSecret))
}
