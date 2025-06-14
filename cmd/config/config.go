package config

import "fmt"

// server params
const (
	Port = ":8080"
)

// urls
const (
	BaseURL = "/api/go-manage-hex"
)

// service errors
var (
	ErrUserNotFound      = fmt.Errorf("user already exists")
	ErrInvalidEmail      = fmt.Errorf("invalid email address")
	ErrInvalidPassword   = fmt.Errorf("invalid password")
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
)
