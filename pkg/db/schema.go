package db

import (
	"errors"
	"time"
)

const (
	AccessLevelBasic = "basic"
	AccessLevelAdmin = "admin"
)

type User struct {
	ID            string     `json:"id"`
	FirstName     string     `json:"firstName"`
	LastName      string     `json:"lastName"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	DOB           *time.Time `json:"dob"`
	VerifiedEmail bool       `json:"verifiedEmail"`
	AccessLevel   string     `json:"accessLevel"`
}

var (
	ErrNotFound = errors.New("not found")
)
