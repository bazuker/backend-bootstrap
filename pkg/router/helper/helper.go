package helper

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	ContextDatabase = "db"
	ContextCache    = "cache"
	ContextUserID   = "userID"
)

type HTTPError struct {
	Message string `json:"message"`
}

type SessionData struct {
	UserID string
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
