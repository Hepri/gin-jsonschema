package schema

import (
	"net/http"

	"io"

	"encoding/json"

	"fmt"

	"sync"

	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
)

var (
	invalidJSONBodyResponse = gin.H{"message": "Invalid json body"}

	cache   map[string]*gojsonschema.Schema = make(map[string]*gojsonschema.Schema)
	cacheMx sync.Mutex
)

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

func buildSchemaFromString(str string) (*gojsonschema.Schema, error) {
	// try to get value from cache without locking, cache filled only once, so we don't need lock for each call
	if sch, found := cache[str]; !found {
		// if value not found, we should create new schema and put it in cache

		// acquire cache lock
		cacheMx.Lock()
		defer cacheMx.Unlock()

		// now read again, probably other goroutine already write value in cache
		if sch, found = cache[str]; found {
			return sch, nil
		}

		// create new schema
		sch, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(str))
		if err != nil {
			return nil, NewErrCannotBuildSchema(err)
		}

		cache[str] = sch
		return sch, nil
	} else {
		return sch, nil
	}
}

func Validate(handler gin.HandlerFunc, schemaStr string) gin.HandlerFunc {
	sch, err := buildSchemaFromString(schemaStr)
	if err != nil {
		panic(fmt.Sprintf("Cannot build schema from string %v", schemaStr))
	}

	return ValidateSchema(handler, sch)
}

func ValidateSchema(handler gin.HandlerFunc, schema *gojsonschema.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateBodyUsingSchema(c.Request, schema); err == nil {
			handler(c)
		} else {
			handleError(c, err)
		}
	}
}

func BindJSON(c *gin.Context, schemaStr string, obj interface{}) error {
	sch, err := buildSchemaFromString(schemaStr)
	if err != nil {
		panic(err)
	}

	return BindJSONSchema(c, sch, obj)
}

func BindJSONSchema(c *gin.Context, schema *gojsonschema.Schema, obj interface{}) (err error) {
	defer func() {
		if err != nil {
			handleError(c, err)
		}
	}()

	// validate body
	if err = validateBodyUsingSchema(c.Request, schema); err != nil {
		return
	}

	// read body and unmarshal json
	var body []byte
	if body, err = ioutil.ReadAll(c.Request.Body); err != nil {
		return
	}

	if err = json.Unmarshal(body, obj); err != nil {
		return
	}

	return
}

func handleError(c *gin.Context, err error) {
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
