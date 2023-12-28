package db

import (
	"errors"
	"time"
)

type User struct {
	ID            string     `json:"id"`
	FirstName     string     `json:"firstName"`
	LastName      string     `json:"lastName"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	DOB           *time.Time `json:"dob"`
	VerifiedEmail bool       `json:"verifiedEmail"`
}

var (
	ErrNotFound = errors.New("not found")
)
