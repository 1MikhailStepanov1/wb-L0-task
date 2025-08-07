package errors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = NewHTTPError(404, "{entity} not found")
)

type HTTPError struct {
	Code     int
	Template string
	entity   string
}

func NewHTTPError(code int, template string) *HTTPError {
	return &HTTPError{Code: code, Template: template}
}

func (e *HTTPError) ForEntity(entity string) *HTTPError {
	return &HTTPError{
		Code:     e.Code,
		Template: e.Template,
		entity:   strings.ToLower(entity),
	}
}

func (e *HTTPError) Error() string {
	msg := e.Template
	if e.entity != "" {
		msg = strings.ReplaceAll(msg, "{entity}", e.entity)
	}
	return fmt.Sprintf("[%d] %s", e.Code, msg)
}

func (e *HTTPError) Is(target error) bool {
	var t *HTTPError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return e.Code == t.Code
}
