package api_errors

import "net/http"

var (
	ErrInternalServerError = "10000"
	ErrUnauthorizedAccess  = "10001"
	ErrInvalidUserID       = "10002"
	ErrValidation          = "10003"
	ErrDeleteFailed        = "10004"

	ErrUserNotFound       = "20005"
	ErrEmailAlreadyExists = "20000"
	ErrEmailNotFound      = "20001"
	ErrInvalidPassword    = "20002"
)

type MessageAndStatus struct {
	Message string
	Status  int
}

var MapErrorCodeMessage = map[string]MessageAndStatus{
	// 10000 - 19999: Common errors
	ErrInternalServerError: {"Internal Server Error", http.StatusInternalServerError},
	ErrUnauthorizedAccess:  {"Unauthorized Access", http.StatusUnauthorized},
	ErrInvalidUserID:       {"Invalid User ID", http.StatusBadRequest},
	ErrValidation:          {"Validation Error", http.StatusBadRequest},
	ErrDeleteFailed:        {"Delete Failed", http.StatusInternalServerError},
	ErrUserNotFound:        {"User Not Found", http.StatusNotFound},
	// Service errors
	ErrEmailAlreadyExists: {"Email Already Exists", http.StatusBadRequest},
	ErrEmailNotFound:      {"Email Not Found", http.StatusBadRequest},
	ErrInvalidPassword:    {"Invalid Password", http.StatusBadRequest},
}
