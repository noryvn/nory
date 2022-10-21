package response

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

func (r *Response[T]) Respond(c *fiber.Ctx) error {
	return c.Status(r.Code).JSON(r)
}

type ResponseError Response[*struct{}]

func (r *ResponseError) Error() string {
	return r.Message
}

func (r *ResponseError) Respond(c *fiber.Ctx) error {
	return c.Status(r.Code).JSON(r)
}

func NewBadRequest(msg string, errs map[string][]string) *ResponseError {
	return &ResponseError{
		Code:    400,
		Message: msg,
		Errors:  errs,
	}
}

func NewUnathorized(msg string) *ResponseError {
	return &ResponseError{
		Code:    401,
		Message: msg,
	}
}

func NewTooManyRequest(msg string) *ResponseError {
	return &ResponseError{
		Code:    429,
		Message: msg,
	}
}

func New[T any](code int, data T) *Response[T] {
	return &Response[T]{
		Code: code,
		Data: data,
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	res, ok := err.(*ResponseError)
	if ok {
		return res.Respond(c)
	}
	return err
}
