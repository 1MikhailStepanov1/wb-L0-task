package errors

import (
	"errors"
	"strings"
)

var (
	ErrNotFound      = NewErrorForEntity(404, "{entity} not found")
	ErrInvalidEntity = NewErrorForEntity(400, "Failed to pass validation field: {entity}")
)

type ErrorForEntity struct {
	Code     int32
	Template string
	entity   string
}

func NewErrorForEntity(code int32, template string) *ErrorForEntity {
	return &ErrorForEntity{Code: code, Template: template}
}

func (e *ErrorForEntity) ForEntity(entity string) *ErrorForEntity {
	return &ErrorForEntity{
		Code:     e.Code,
		Template: e.Template,
		entity:   strings.ToLower(entity),
	}
}

func (e *ErrorForEntity) Error() string {
	msg := e.Template
	if e.entity != "" {
		msg = strings.ReplaceAll(msg, "{entity}", e.entity)
	}
	return msg
}

func (e *ErrorForEntity) Is(target error) bool {
	var t *ErrorForEntity
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return e.Code == t.Code
}
