package schema

import (
	"fmt"
	"strings"
)

type ErrCannotBuildSchema struct {
	err error
}

func (e *ErrCannotBuildSchema) Error() string {
	return fmt.Sprintf("Cannot build schema: %v", e.err)
}

func NewErrCannotBuildSchema(err error) *ErrCannotBuildSchema {
	return &ErrCannotBuildSchema{err}
}

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
