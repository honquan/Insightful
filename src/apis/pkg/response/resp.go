package response

import (
	"fmt"
	"insightful/src/apis/pkg/pagination"
	"net/http"
)

type Response interface {
	GetData() interface{}
	GetCode() int
}

type ValidateErrMeta struct {
	Key  string `json:"_key"`
	Name string `json:"name"`
}

type Meta struct {
	*pagination.Pagination `json:",omitempty"`
	*ValidateErrMeta       `json:",omitempty"`
	Code                   int         `json:"code"`
	Message                string      `json:"message"`
	Debug                  interface{} `json:"debug,omitempty"`
}

type response struct {
	Code int         `json:"-"`
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

func (r *response) GetData() interface{} {
	return r.Data
}

func (r *response) GetCode() int {
	return r.Code
}

func SuccessPagination(data interface{}, p pagination.Pagination, msg string, args ...interface{}) Response {
	return &response{
		Code: http.StatusOK,
		Meta: Meta{
			Pagination: &p,
			Code:       http.StatusOK,
			Message:    fmt.Sprintf(msg, args...),
		},
		Data: data,
	}
}

func Info(template string, args ...interface{}) Response {
	return Success(nil, template, args...)
}

func Success(data interface{}, msg string, arg ...interface{}) Response {
	return &response{
		Code: http.StatusOK,
		Meta: Meta{
			Code:    http.StatusOK,
			Message: fmt.Sprintf(msg, arg...),
		},
		Data: data,
	}
}

func Error(code int, msg string, arg ...interface{}) Response {
	return &response{
		Code: code,
		Meta: Meta{
			Code:    code,
			Message: fmt.Sprintf(msg, arg...),
		},
	}
}

func ErrorDebug(code int, msg string, debug interface{}) Response {
	return &response{
		Code: code,
		Meta: Meta{
			Code:    code,
			Message: msg,
			Debug:   debug,
		},
	}
}

func BadRequest(msg string, args ...interface{}) Response {
	return &response{
		Code: http.StatusBadRequest,
		Meta: Meta{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf(msg, args...),
		},
	}
}

func InternalError(err error) Response {
	return &response{
		Code: http.StatusInternalServerError,
		Meta: Meta{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}
