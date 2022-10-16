package interfaces

type Response[T any] struct {
	Code    int
	Data    T                   `json:",omitempty"`
	Message string              `json:",omitempty"`
	Errors  map[string][]string `json:",omitempty`
}

type ResponsePaginate[T any] struct {
	Code int
	Data []T

	Pagination struct {
		Page       int
		Items      int
		TotalItems int
	}
}

type ResponseError Response[*struct{}]

func (r ResponseError) Error() string {
	return r.Message
}

func NewResponseUnathorized(msg string) ResponseError {
	return ResponseError{
		Code:    401,
		Message: msg,
	}
}

func NewResponse[T any](code int, data T) Response[T] {
	return Response[T]{
		Code: code,
		Data: data,
	}
}
