package helper

import (
	"crypto/rand"
	"encoding/base64"
)

// Context keys that can be used to fetch useful information
// from a request's context.
// Context availability depends on the middleware called before a handler.
const (
	ContextDatabase        = "db"
	ContextCache           = "cache"
	ContextUserID          = "userID"
	ContextUserAccessLevel = "userAccessLevel"
)

type HTTPError struct {
	Message string `json:"message"`
}

type SessionData struct {
	UserID      string
	AccessLevel string
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
