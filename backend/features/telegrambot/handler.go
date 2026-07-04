package telegrambot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"tma-boilerplate/config"
	"tma-boilerplate/database"

	"github.com/gin-gonic/gin"
)

// ---- Telegram API types ----

type TelegramUpdate struct {
	UpdateID int              `json:"update_id"`
	Message  *TelegramMessage `json:"message"`
}

type TelegramMessage struct {
	MessageID int          `json:"message_id"`
	From      TelegramUser `json:"from"`
	Chat      TelegramChat `json:"chat"`
	Text      string       `json:"text"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type TelegramChat struct {
	ID int64 `json:"id"`
}

// ---- Webhook handler ----

func HandleWebhook(c *gin.Context) {
	var update TelegramUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		c.Status(http.StatusOK)
		return
	}

	if update.Message != nil {
		handleMessage(update.Message)
	}

	c.Status(http.StatusOK)
}

func handleMessage(msg *TelegramMessage) {
	// A guest who tapped the "Continue in Telegram" link on the website
	// arrives here as "/start <guest_id>" — link their web progress to
	// this Telegram account before replying.
	if guestID, ok := startPayload(msg.Text); ok {
		linkGuestAccount(guestID, msg.From)
	}
	sendWelcome(msg.Chat.ID, msg.From.FirstName)
}

// startPayload extracts the deep-link payload from a "/start <payload>" command.
func startPayload(text string) (string, bool) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "/start") {
		return "", false
	}
	payload := strings.TrimSpace(strings.TrimPrefix(text, "/start"))
	if payload == "" || len(payload) > 64 {
		return "", false
	}
	return payload, true
}

// linkGuestAccount attaches a Telegram identity to a guest account created on
// the website, so the guest's progress carries over instead of starting a
// brand-new account. It's a no-op if the guest is already linked or unknown.
func linkGuestAccount(guestID string, from TelegramUser) {
	_, err := database.DB.Exec(`
		UPDATE users SET telegram_id = $1, username = $2, first_name = $3, updated_at = NOW()
		WHERE guest_id = $4 AND telegram_id IS NULL`,
		from.ID, from.Username, from.FirstName, guestID,
	)
	if err != nil {
		log.Printf("linkGuestAccount: failed to link guest %s: %v", guestID, err)
	}
}

// ---- Message builders ----

func sendWelcome(chatID int64, firstName string) {
	text := fmt.Sprintf("Hello, %s! Welcome.", firstName)

	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]interface{}{
			{
				{
					"text": "Open App",
					"web_app": map[string]string{
						"url": config.App.TelegramMiniAppURL,
					},
				},
			},
		},
	}

	payload := map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": keyboard,
	}

	if err := telegramAPICall("sendMessage", payload); err != nil {
		log.Printf("sendMessage error: %v", err)
	}
}

// ---- Webhook registration ----

type ngrokTunnelsResponse struct {
	Tunnels []struct {
		PublicURL string `json:"public_url"`
		Proto     string `json:"proto"`
	} `json:"tunnels"`
}

func discoverNgrokURL(apiURL string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	for i := 0; i < 20; i++ {
		resp, err := client.Get(apiURL + "/api/tunnels")
		if err == nil {
			var data ngrokTunnelsResponse
			json.NewDecoder(resp.Body).Decode(&data)
			resp.Body.Close()
			for _, t := range data.Tunnels {
				if t.Proto == "https" {
					return t.PublicURL, nil
				}
			}
		}
		log.Printf("Waiting for ngrok tunnel... (%d/20)", i+1)
		time.Sleep(2 * time.Second)
	}
	return "", fmt.Errorf("ngrok tunnel not available after 40s")
}

func RegisterWebhook() error {
	if config.App.TelegramBotToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is not set")
	}

	baseURL := config.App.WebhookURL
	if baseURL == "" && config.App.NgrokAPIURL != "" {
		discovered, err := discoverNgrokURL(config.App.NgrokAPIURL)
		if err != nil {
			return err
		}
		baseURL = discovered
		log.Printf("Discovered ngrok URL: %s", baseURL)
		// This branch only runs in dev (ngrok is only wired up via docker-compose.dev.yml),
		// so always point the Mini App at the freshly discovered tunnel, overriding
		// whatever placeholder is left in TELEGRAM_MINI_APP_URL.
		config.App.TelegramMiniAppURL = baseURL
		log.Printf("Mini App URL set to ngrok: %s", baseURL)
	}
	if baseURL == "" {
		return fmt.Errorf("WEBHOOK_URL is not set (and NGROK_API_URL is not configured)")
	}

	webhookURL := baseURL + "/api/webhook/telegram"

	payload := map[string]interface{}{
		"url":                  webhookURL,
		"allowed_updates":      []string{"message"},
		"drop_pending_updates": true,
	}

	if err := telegramAPICall("setWebhook", payload); err != nil {
		return fmt.Errorf("setWebhook failed: %w", err)
	}

	log.Printf("Telegram webhook registered: %s", webhookURL)
	return nil
}

// GetBotInfo exposes the bot username so the frontend can build a
// "https://t.me/<username>" deep link — no auth needed, it's public info.
func GetBotInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"username": config.App.TelegramBotUsername})
}

// GetWebhookInfo returns current webhook status.
func GetWebhookInfo(c *gin.Context) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getWebhookInfo", config.App.TelegramBotToken)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var result interface{}
	json.Unmarshal(body, &result)
	c.JSON(http.StatusOK, result)
}

// ---- Helper ----

func telegramAPICall(method string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", config.App.TelegramBotToken, method)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if !result.OK {
		return fmt.Errorf("telegram API error: %s", result.Description)
	}
	return nil
}
