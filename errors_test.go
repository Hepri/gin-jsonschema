package schema

import (
	"errors"
	"testing"
)

func TestNewErrSchemaValidationEmpty(t *testing.T) {
	er := NewErrSchemaValidation(nil)

	// test empty string
	if str := er.Error(); str != "" {
		t.Errorf("invalid error string, received '%v' expected ''", str)
	}
}

func TestNewErrSchemaValidationTwoErrors(t *testing.T) {
	var arr = []string{
		"error 1",
		"error 2",
	}
	var errorsStr string = "error 1; error 2"

	er := NewErrSchemaValidation(arr)

	// test empty string
	if str := er.Error(); str != errorsStr {
		t.Errorf("invalid error string, received '%v' expected '%v'", str, errorsStr)
	}
}

func TestNewErrCannotBuildSchema(t *testing.T) {
	er := NewErrCannotBuildSchema(errors.New("test str"))
	refVal := `Cannot build schema: test str`

	if str := er.Error(); str != refVal {
		t.Errorf("invalid error string, received '%v' expected '%v'", str, refVal)
	}
}
