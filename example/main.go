package main

import (
	"github.com/gin-gonic/gin"

	"net/http"

	"github.com/Hepri/gin-jsonschema"
)

var testSchema string = `
{
    "title": "Test Schema",
    "type": "object",
    "properties": {
        "value": {
            "type": "integer"
        }
    },
    "required": ["value"]
}
`

type testBody struct {
	Value int `json:"value"`
}

func handlerFunc(c *gin.Context) {
	c.Status(http.StatusOK)
}

func bindJSONHandler(c *gin.Context) {
	var js testBody
	if schema.BindJSON(c, testSchema, &js) == nil {
		c.JSON(http.StatusOK, gin.H{
			"value": js.Value,
		})
	}
}

func main() {
	r := gin.Default()

	// wrap handler, all invalid json schema requests will produce Bad Request
	r.POST("/validate", schema.Validate(handlerFunc, testSchema))

	// use BindJSON inside handler
	r.POST("/bind", bindJSONHandler)
	r.Run(":8080")
}
