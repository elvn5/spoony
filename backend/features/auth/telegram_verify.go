package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"tma-boilerplate/config"
)

// VerifyTelegramInitData validates the Telegram Web App initData signature.
// Returns parsed query values on success.
func VerifyTelegramInitData(initData string) (url.Values, error) {
	vals, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("invalid initData format")
	}

	hash := vals.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("missing hash in initData")
	}

	// Validate timestamp (allow 24h window)
	authDate := vals.Get("auth_date")
	if authDate == "" {
		return nil, fmt.Errorf("missing auth_date in initData")
	}
	var ts int64
	if _, err := fmt.Sscanf(authDate, "%d", &ts); err != nil {
		return nil, fmt.Errorf("invalid auth_date")
	}
	if time.Now().Unix()-ts > 86400 {
		return nil, fmt.Errorf("initData expired")
	}

	// Build the data-check string
	var pairs []string
	for k, v := range vals {
		if k == "hash" {
			continue
		}
		pairs = append(pairs, k+"="+v[0])
	}
	sort.Strings(pairs)
	dataCheckString := strings.Join(pairs, "\n")

	// Derive secret key: HMAC-SHA256("WebAppData", botToken)
	secretKey := hmacSHA256([]byte("WebAppData"), []byte(config.App.TelegramBotToken))
	expectedHash := hex.EncodeToString(hmacSHA256(secretKey, []byte(dataCheckString)))

	if !hmac.Equal([]byte(expectedHash), []byte(hash)) {
		return nil, fmt.Errorf("invalid signature")
	}

	return vals, nil
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
