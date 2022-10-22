package response

import "github.com/gofiber/fiber/v2"

type Response[T any] struct {
	Code    int    `json:"code"`
	Data    T      `json:"data"`
	Message string `json:"message,omitempty"`
}

func (r *Response[T]) Respond(c *fiber.Ctx) error {
	if r.Code == 204 {
		return c.SendStatus(204)
	}
	return c.Status(r.Code).JSON(r)
}

func New[T any](code int, data T) *Response[T] {
	return &Response[T]{
		Code: code,
		Data: data,
	}
}

type ResponseError Response[*struct{}]

func (r *ResponseError) Error() string {
	return r.Message
}

func (r *ResponseError) Respond(c *fiber.Ctx) error {
	return c.Status(r.Code).JSON(r)
}

func NewBadRequest(msg string) *ResponseError {
	return &ResponseError{
		Code:    400,
		Message: msg,
	}
}

func NewUnathorized(msg string) *ResponseError {
	return &ResponseError{
		Code:    401,
		Message: msg,
	}
}

func NewNotFound(msg string) *ResponseError {
	return &ResponseError{
		Code:    404,
		Message: msg,
	}
}

func NewTooManyRequest(msg string) *ResponseError {
	return &ResponseError{
		Code:    429,
		Message: msg,
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	res, ok := err.(*ResponseError)
	if ok {
		return res.Respond(c)
	}
	return err
}
