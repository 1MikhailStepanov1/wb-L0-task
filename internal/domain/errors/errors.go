package errors

import (
	"errors"
	"strings"
)

var (
	ErrNotFound      = NewEntityError(404, "{entity} not found")
	ErrInvalidEntity = NewEntityError(400, "Failed to pass validation field: {entity}")
	ErrBrokenEntity  = NewEntityError(400, "Invalid entity received: {entity}")
)

type EntityError struct {
	Code     int32
	Template string
	entity   string
}

func NewEntityError(code int32, template string) *EntityError {
	return &EntityError{Code: code, Template: template}
}

func (e *EntityError) ForEntity(entity string) *EntityError {
	return &EntityError{
		Code:     e.Code,
		Template: e.Template,
		entity:   strings.ToLower(entity),
	}
}

func (e *EntityError) Error() string {
	msg := e.Template
	if e.entity != "" {
		msg = strings.ReplaceAll(msg, "{entity}", e.entity)
	}
	return msg
}

func (e *EntityError) Is(target error) bool {
	var t *EntityError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return e.Code == t.Code
}
