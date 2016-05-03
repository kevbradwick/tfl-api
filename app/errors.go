package app

import "fmt"

type HttpError struct {
	Code    int
	Message string
}

func (h *HttpError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", h.Code, h.Message)
}

type DbError struct {
	Message string
	Query   interface{}
}

func (d *DbError) Error() string {
	return fmt.Sprintf("A Mongo query error occured. %q", d.Query)
}
