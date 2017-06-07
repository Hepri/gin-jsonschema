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

// validate using schema as gojsonschema.JSONLoader
loader := gojsonschema.NewStringLoader()
r.POST("/loader", schema.ValidateJSONLoader(handlerFunc, loader))

// validate using schema as *gojsonschema.Schema
sc, _ := gojsonschema.NewSchema(loader)
r.POST("/schema", schema.Validate(handlerFunc, sc))
```

In order to use json schema from file or other sources look at [gojsonschema documentation](https://github.com/xeipuuv/gojsonschema)


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
