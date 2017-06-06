package schema

import "testing"

func TestNewErrSchemaValidationEmpty(t *testing.T) {
	er := NewErrSchemaValidation(nil)

	// test empty string
	if str := er.Error(); str != "" {
		t.Errorf("invalid error string, received '%v' expected ''", str)
	}
}

func TestNewErrSchemaValidationTwoErrors(t *testing.T) {
	var errors = []string{
		"error 1",
		"error 2",
	}
	var errorsStr string = "error 1; error 2"

	er := NewErrSchemaValidation(errors)

	// test empty string
	if str := er.Error(); str != errorsStr {
		t.Errorf("invalid error string, received '%v' expected '%v'", str, errorsStr)
	}
}
