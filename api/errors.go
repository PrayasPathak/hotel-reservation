package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const (
	ErrTokenExpiredMsg = "token expired"
	ErrBadRequestMsg   = "invalid JSON request"
	ErrInvalidIDMsg    = "invalid id"
	ErrUnauthorizedMsg = "unauthorized"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func (e Error) Error() string {
	return e.Err
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  ErrInvalidIDMsg,
	}
}

func ErrUnauthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  ErrUnauthorizedMsg,
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  ErrBadRequestMsg,
	}
}

func ErrResourceNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  fmt.Sprintf("%s resource not found", res),
	}
}
