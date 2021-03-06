package schema

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"encoding/json"

	"strings"

	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
)

var testSchema = `
{
    "title": "test",
    "type": "object",
    "properties": {
        "value1": {
            "type": "integer"
        },
        "value2": {
            "type": "boolean"
        }
    },
    "required": ["value1"],
    "additionalProperties": false
}`

type testBody struct {
	Value1 int  `json:"value1"`
	Value2 bool `json:"value2"`
}

func TestBuildSchemaFromString(t *testing.T) {
	sch, err := buildSchemaFromString(testSchema)
	if err != nil {
		t.Fatalf("cannot build schema from string: %v", err)
	}

	// call again with same schema
	sch2, err2 := buildSchemaFromString(testSchema)
	if err2 != nil {
		t.Fatalf("cannot build schema from string: %v", err)
	}

	if sch != sch2 {
		t.Errorf("expected schema to be equal between calls: %v", err)
	}

	// build invalid schema error
	_, err = buildSchemaFromString(`{"a": `)
	if err == nil {
		t.Errorf("build should return error, but received nil")
	}
}

func TestBuildSchemaFromStringInvalidSchema(t *testing.T) {
	_, err := buildSchemaFromString(`{"a"`)
	if err == nil {
		t.Errorf("build schema should return error, but returns nil")
	} else {
		switch err.(type) {
		case *ErrCannotBuildSchema:
		default:
			t.Errorf("build schema should return ErrCannotBuildSchema, but returns %v", err)
		}
	}
}

func createRequestWithBody(t *testing.T, body io.Reader) *http.Request {
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatalf("cannot create new request: %v", err)
	}

	return req
}

func TestValidateBodyUsingSchemaNilBody(t *testing.T) {
	req := createRequestWithBody(t, nil)
	err := validateBodyUsingSchema(req, nil)
	if err != io.EOF {
		t.Errorf("validate returns unexpected error, want %v received %v", io.EOF, err)
		return
	}
}

func TestValidateBodyUsingSchemaEmptyBody(t *testing.T) {
	req := createRequestWithBody(t, bytes.NewReader([]byte("")))
	err := validateBodyUsingSchema(req, nil)
	if err != io.EOF {
		t.Errorf("validate returns unexpected error, want %v received %v", io.EOF, err)
		return
	}
}

func createSchemaFromString(t *testing.T, str string) *gojsonschema.Schema {
	sch, err := buildSchemaFromString(str)
	if err != nil {
		t.Fatalf("cannot create schema: %v", err)
	}

	return sch
}

func TestValidateBodyUsingSchemaNonJSONBody(t *testing.T) {
	ch := createSchemaFromString(t, `{}`)

	// check json syntax error
	req := createRequestWithBody(t, strings.NewReader(`{a: 1`))
	err := validateBodyUsingSchema(req, ch)
	if err == nil {
		t.Errorf("validate should return syntax error, but nil received")
	} else {
		switch v := err.(type) {
		case *json.SyntaxError:
			// ok
		default:
			t.Errorf("validate should return %T received %T:%v", &json.SyntaxError{}, v, err)
		}
	}

	req2 := createRequestWithBody(t, strings.NewReader(`{"a": 1`))
	err = validateBodyUsingSchema(req2, ch)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("validate should return %v received %v", io.ErrUnexpectedEOF, err)
	}
}

func TestValidateBodyUsingSchemaInvalidJSON(t *testing.T) {
	sc := createSchemaFromString(t, testSchema)

	// bad value1 type
	req := createRequestWithBody(t, strings.NewReader(`{"value1": "bad type"}`))
	err := validateBodyUsingSchema(req, sc)
	if err == nil {
		t.Errorf("validate should return error, but received nil")
	}

	// extra field that not present in schema
	req2 := createRequestWithBody(t, strings.NewReader(`{"value1": 1, "value3": 3}`))
	err = validateBodyUsingSchema(req2, sc)
	if err == nil {
		t.Errorf("validate should return error, but received nil")
	}

	// valid json
	req3 := createRequestWithBody(t, strings.NewReader(`{"value1": 1, "value2": true}`))
	err = validateBodyUsingSchema(req3, sc)
	if err != nil {
		t.Errorf("validate should return nil, but received %v", err)
	}
}

func getTestServer(handlerFunc gin.HandlerFunc) *httptest.Server {
	// disable debug message in console
	gin.SetMode(gin.ReleaseMode)

	// using `New` instead of `Default` to disable logging and recovery
	r := gin.New()
	r.POST("/", handlerFunc)

	return httptest.NewServer(r)
}

func okHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}

func testRequest(t *testing.T, url string, body io.Reader, expectedStatusCode int) {
	res, err := http.Post(url, "application/json", body)
	if err != nil {
		t.Fatalf("error happens during post request %v", err)
	}

	if res.StatusCode != expectedStatusCode {
		t.Errorf("expected response status code %v, but received %v", expectedStatusCode, res.StatusCode)
	}
}

func testHandlerResponses(t *testing.T, handler gin.HandlerFunc) {
	// ensure server validate requests using `testSchema`
	ts := getTestServer(handler)
	defer ts.Close()

	// nil body
	testRequest(t, ts.URL, nil, http.StatusBadRequest)
	// empty body
	testRequest(t, ts.URL, strings.NewReader(``), http.StatusBadRequest)
	// non json, invalid syntax
	testRequest(t, ts.URL, strings.NewReader(`{a: 1`), http.StatusBadRequest)
	// non json, unexpected eof
	testRequest(t, ts.URL, strings.NewReader(`{"a": 1`), http.StatusBadRequest)
	// non json, unexpected eof
	testRequest(t, ts.URL, strings.NewReader(`{"a": 1`), http.StatusBadRequest)

	// invalid json, wrong value1 type
	testRequest(t, ts.URL, strings.NewReader(`{"value1": "aaa"}`), http.StatusBadRequest)
	// invalid json, additional field that not present in schema
	testRequest(t, ts.URL, strings.NewReader(`{"value1": 1, "value3": "aaa"}`), http.StatusBadRequest)
	// invalid json, missing required field
	testRequest(t, ts.URL, strings.NewReader(`{"value2": true}`), http.StatusBadRequest)

	// valid json
	testRequest(t, ts.URL, strings.NewReader(`{"value1": 1, "value2": true}`), http.StatusOK)
}

func TestValidate(t *testing.T) {
	testHandlerResponses(t, Validate(okHandler, testSchema))
}

func TestValidateSchema(t *testing.T) {
	// build schema from string
	sc, err := buildSchemaFromString(testSchema)
	if err != nil {
		t.Errorf("cannot build schema from string %v", testSchema)
	}
	testHandlerResponses(t, ValidateSchema(okHandler, sc))
}

func makeBindJSONHandler(t *testing.T) gin.HandlerFunc {
	return func(c *gin.Context) {
		var js testBody
		BindJSON(c, testSchema, &js)
		c.Status(http.StatusOK)
	}
}

func makeBindJSONSchemaHandler(t *testing.T) gin.HandlerFunc {
	sch := createSchemaFromString(t, testSchema)

	return func(c *gin.Context) {
		var js testBody
		BindJSONSchema(c, sch, &js)
		c.Status(http.StatusOK)
	}
}

func TestBindJSONValidate(t *testing.T) {
	testHandlerResponses(t, Validate(makeBindJSONHandler(t), testSchema))
}

func TestBindJSONSchemaValidate(t *testing.T) {
	testHandlerResponses(t, Validate(makeBindJSONSchemaHandler(t), testSchema))
}

func makeBindJSONUnmarshalHandler(t *testing.T) gin.HandlerFunc {
	return func(c *gin.Context) {
		var js testBody
		BindJSON(c, testSchema, &js)
		if js.Value1 != 2 {
			t.Errorf("value1 expected to be %v but received %v", 2, js.Value1)
		}
		if js.Value2 != true {
			t.Errorf("value2 expected to be %v but received %v", true, js.Value2)
		}
		c.Status(http.StatusOK)
	}
}

func makeBindJSONSchemaUnmarshalHandler(t *testing.T) gin.HandlerFunc {
	sch := createSchemaFromString(t, testSchema)

	return func(c *gin.Context) {
		var js testBody
		BindJSONSchema(c, sch, &js)
		if js.Value1 != 2 {
			t.Errorf("value1 expected to be %v but received %v", 2, js.Value1)
		}
		if js.Value2 != true {
			t.Errorf("value2 expected to be %v but received %v", true, js.Value2)
		}
		c.Status(http.StatusOK)
	}
}

func testUnmarshalHandler(t *testing.T, handler gin.HandlerFunc) {
	// ensure server validate requests using `testSchema`
	ts := getTestServer(handler)
	defer ts.Close()

	testRequest(t, ts.URL, strings.NewReader(`{"value1": 2, "value2": true}`), http.StatusOK)
}

func TestBindJSONUnmarshal(t *testing.T) {
	testUnmarshalHandler(t, Validate(makeBindJSONUnmarshalHandler(t), testSchema))
}

func TestBindJSONSchemaUnmarshal(t *testing.T) {
	testUnmarshalHandler(t, Validate(makeBindJSONSchemaUnmarshalHandler(t), testSchema))
}
