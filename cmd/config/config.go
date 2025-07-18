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
	ErrSearchingUser = "error searching user"
	ErrCreatingUser  = "error creating user"
	ErrDeletingUser  = "error deleting user"
	ErrUpdatingUser  = "error updating user"
	ErrChangingPwd   = "error changing password"

	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrInvalidEmail      = fmt.Errorf("invalid email address")
	ErrInvalidPassword   = fmt.Errorf("invalid password")
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
)

//handler messages

const (
	UserCreatedMsg           = "user created successfully"
	UserFoundMsg             = "user found successfully"
	InvalidBodyMsg           = "invalid body params"
	InvalidQueryParamsMsg    = "invalid query params"
	InvalidConfirmationMsg   = "delete operation canceled"
	InvalidConfirmationValue = "invalid confirmation value"
	DeleteCancelledMsg       = "delete operation canceled"
	UserDeletedMsg           = "user deleted successfully"
	UserUpdatedMsg           = "user updated succesfully"
	UserPwdChangeMsg         = "password changed successfully"
)
