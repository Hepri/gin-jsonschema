package schema

import (
	"net/http"

	"io"

	"encoding/json"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
)

var invalidJSONBodyResponse = gin.H{"message": "Invalid json body"}

func validateBodyUsingSchema(req *http.Request, schema *gojsonschema.Schema) error {
	body, err := drainHTTPRequestBody(req)
	if err != nil {
		return err
	}

	// validate body
	result, err := schema.Validate(gojsonschema.NewBytesLoader(body))
	if err != nil {
		return err
	}

	// schema not valid, create validation error
	if !result.Valid() {
		var errors []string
		for _, er := range result.Errors() {
			errors = append(errors, er.String())
		}
		return NewErrSchemaValidation(errors)
	}

	return nil
}

func Validate(handler gin.HandlerFunc, schema *gojsonschema.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateBodyUsingSchema(c.Request, schema); err == nil {
			handler(c)
		} else {
			c.Abort()
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				c.JSON(http.StatusBadRequest, invalidJSONBodyResponse)
			} else {
				switch v := err.(type) {
				case *json.SyntaxError:
					c.JSON(http.StatusBadRequest, invalidJSONBodyResponse)
				case *ErrSchemaValidation:
					c.JSON(http.StatusBadRequest, gin.H{
						"messages": v.Errors,
					})
				default:
					c.Status(http.StatusInternalServerError)
				}
			}
		}
	}
}

func ValidateString(handler gin.HandlerFunc, str string) gin.HandlerFunc {
	loader := gojsonschema.NewStringLoader(str)
	sch, err := gojsonschema.NewSchema(loader)
	if err != nil {
		panic(fmt.Sprintf("Cannot build schema from string %v", str))
	}

	return Validate(handler, sch)
}

func ValidateJSONLoader(handler gin.HandlerFunc, loader gojsonschema.JSONLoader) gin.HandlerFunc {
	sch, err := gojsonschema.NewSchema(loader)
	if err != nil {
		panic(fmt.Sprintf("Cannot build schema from loader %v", loader))
	}

	return Validate(handler, sch)
}
