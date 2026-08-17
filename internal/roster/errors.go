package roster

import (
	"errors"
	"fmt"
)

var ErrInvalidSchedule = errors.New("invalid schedule")

type FieldError struct {
	Field   string
	Message string
	Err     error
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (e *FieldError) Unwrap() error {
	return nil
}

func invalidField(field, message string) error {
	return &FieldError{Field: field, Message: message, Err: ErrInvalidSchedule}
}
