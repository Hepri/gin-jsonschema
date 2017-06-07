# gin-jsonschema: JSON Schema validation wrapper for Gin [![Build Status](https://travis-ci.org/Hepri/case-transformer.png?branch=master)](https://travis-ci.org/Hepri/gin-jsonschema)

## Installation

```
go get github.com/Hepri/gin-jsonschema
```
   
Dependencies :
* [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
* [github.com/xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema)

## Usage

```
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

type TestBody struct {
    Value int `json:"value"`
}


// using BindJSON
r.POST("/bind", func(c *gin.Context) {
    var js TestBody
    if schema.BindJSON(c, testSchema, &js) == nil {
        // do stuff
    }
})


// validate using schema as string
r.POST("/string", schema.Validate(handlerFunc, testSchema))


// validate using schema as *gojsonschema.Schema
loader := gojsonschema.NewStringLoader(testSchema)
sc, _ := gojsonschema.NewSchema(loader)
r.POST("/schema", schema.ValidateSchema(handlerFunc, sc))


// using wrapper inside handler
var someHandler = func(c *gin.Context) {
    return schema.Validate(func (c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "OK",
        })
    }, testSchema)
}
r.POST("/wrap_inside", someHandler)
```

Read possible ways to build `*gojsonschema.Schema` in [documentation](https://github.com/xeipuuv/gojsonschema)


## Example HTTP Service


```
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

```
