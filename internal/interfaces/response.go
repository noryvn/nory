package interfaces

import "github.com/gofiber/fiber/v2"

type Response[T any] struct {
	Code    int                 `json:"code"`
	Data    T                   `json:"data,omitempty"`
	Message string              `json:"message,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

type ResponsePaginate[T any] struct {
	Code int `json:"code"`
	Data []T `json:"data"`

	Pagination struct {
		Page       int `json:"page"`
		Items      int `json:"items"`
		TotalItems int `json:"totalItems"`
	} `json:"pagination"`
}

type ResponseError = Response[*struct{}]

func (r *Response[T]) Error() string {
	return r.Message
}

func (r *Response[T]) Respond(c *fiber.Ctx) error {
	return c.Status(r.Code).JSON(r)
}

func NewResponseBadRequest(msg string, errs map[string][]string) *ResponseError {
	return &ResponseError{
		Code:    400,
		Message: msg,
		Errors:  errs,
	}
}

func NewResponseUnathorized(msg string) *ResponseError {
	return &ResponseError{
		Code:    401,
		Message: msg,
	}
}

func NewResponseTooManyRequest(msg string) *ResponseError {
	return &ResponseError{
		Code:    429,
		Message: msg,
	}
}

func NewResponse[T any](code int, data T) Response[T] {
	return Response[T]{
		Code: code,
		Data: data,
	}
}
