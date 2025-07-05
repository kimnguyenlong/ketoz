package entity

import "net/http"

const (
	CodeSuccess            = 0
	CodeInternalError      = -500
	CodeInvalidParamsError = -400
)

const (
	MessageSuccess       = "Success"
	MessageInternalError = "Something went wrong!"
)

type Error struct {
	Code     int
	Message  string
	HttpCode int
}

func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	return err.Message
}

func NewInternalError(msg string) *Error {
	return &Error{
		Code:     CodeInternalError,
		Message:  msg,
		HttpCode: http.StatusInternalServerError,
	}
}

func NewInvalidParamsError(msg string) *Error {
	return &Error{
		Code:     CodeInvalidParamsError,
		Message:  msg,
		HttpCode: http.StatusBadRequest,
	}
}
