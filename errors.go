package schema

import (
	"strings"
)

type ErrSchemaValidation struct {
	Errors []string
}

func (e *ErrSchemaValidation) Error() string {
	return strings.Join(e.Errors, "; ")
}

func NewErrSchemaValidation(errors []string) *ErrSchemaValidation {
	return &ErrSchemaValidation{
		Errors: errors,
	}
}
