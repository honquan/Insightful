package errors

import (
	"errors"
	"fmt"
)

type ErrorType uint

const (
	ErrUnknown ErrorType = iota
	ErrBadRequest
	ErrInternal
	ErrNotFound
	ErrForbidden
	ErrUnauthorized
)

type Error struct {
	Type ErrorType
	base error
	msg  string
}

func (e Error) Error() string {
	if e.msg != "" {
		return e.msg
	}
	if e.base != nil {
		return e.base.Error()
	}
	return "Unknown error"
}

func (e Error) Unwrap() error {
	return e.base
}

func BadRequest(msg string, args ...interface{}) error {
	return Error{
		Type: ErrBadRequest,
		base: nil,
		msg:  fmt.Sprintf(msg, args...),
	}
}

func NotFound(msg string, args ...interface{}) error {
	return Error{
		Type: ErrNotFound,
		base: nil,
		msg:  fmt.Sprintf(msg, args...),
	}
}

func Forbidden(msg string, args ...interface{}) error {
	return Error{
		Type: ErrForbidden,
		base: nil,
		msg:  fmt.Sprintf(msg, args...),
	}
}

func Unauthorized(msg string, args ...interface{}) error {
	return Error{
		Type: ErrUnauthorized,
		base: nil,
		msg:  fmt.Sprintf(msg, args...),
	}
}

func InternalError(msg string, args ...interface{}) error {
	return Error{
		Type: ErrInternal,
		base: nil,
		msg:  fmt.Sprintf(msg, args...),
	}
}

func Wrap(err error, t ErrorType) error {
	return Error{
		Type: t,
		base: err,
		msg:  "",
	}
}

func New(template string, args ...interface{}) error {
	return Error{
		Type: ErrUnknown,
		base: nil,
		msg:  fmt.Sprintf(template, args...),
	}
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}
