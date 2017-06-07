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

// validate using schema as string
r.POST("/string", schema.ValidateString(handlerFunc, testSchema))

// validate using schema as *gojsonschema.Schema
loader := gojsonschema.NewStringLoader(testSchema)
sc, _ := gojsonschema.NewSchema(loader)
r.POST("/schema", schema.Validate(handlerFunc, sc))


// using wrapper inside handler
var someHandler = func(c *gin.Context) {
    return schema.ValidateString(func (c *gin.Context) {
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
	r.POST("/", schema.ValidateString(handlerFunc, testSchema))
	r.Run(":8080")
}
```
