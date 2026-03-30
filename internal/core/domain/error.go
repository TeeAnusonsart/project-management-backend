package domain

import "errors"

type ErrorResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var (
    ErrDeviceAlreadyPaired = errors.New("device already paired")
    ErrDeviceNotFound      = errors.New("device not found")
	ErrWidgetNotFound = errors.New("widget not found")
	ErrRoomNotFound   = errors.New("room not found")
	ErrRoomExist = errors.New("room is already exist")
	ErrEmailAlreadyExist = errors.New("Email is already exist")
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidDeviceKey = errors.New("invalid device key")
	ErrPairFailed = errors.New("pairing failed")
)