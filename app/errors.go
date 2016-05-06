package app

type HttpError interface {
	StatusCode() int
	Message() string
}

type GenericHttpError struct {
	code    int
	message string
}

func (h *GenericHttpError) StatusCode() int {
	return h.code
}

func (h *GenericHttpError) Message() string {
	return h.message
}
