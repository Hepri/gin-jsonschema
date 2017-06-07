package main

import (
	"github.com/gin-gonic/gin"

	"net/http"

	"github.com/Hepri/gin-jsonschema"
)

func handlerFunc(c *gin.Context) {
	c.Status(http.StatusOK)
}

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

func main() {
	r := gin.Default()
	// wrap handler, all invalid json schema requests will produce Bad Request
	r.POST("/", schema.Validate(handlerFunc, testSchema))
	r.Run(":8080")
}
