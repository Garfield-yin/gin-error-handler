# Gin Web Framework error handle
A simple error handling middleware for Gin Web Framework.

## Installation
1. Download and install it:

```sh
$ go get -u github.com/garfield-yin/gin-error-handler
```

2. Import it in your code:

```go
import "github.com/garfield-yin/gin-error-handler"
```

## Quick start

```go
package main

import (
  "github.com/gin-gonic/gin"
  "github.com/garfield-yin/gin-error-handler"
  myErrors "github.com/garfield-yin/gin-error-handler/errors"
)

func main() {
  r := gin.Default()
  var errWriter io.Writer = os.Stderr
  /**
   Register your error message flags,like this
   var MsgFlags = map[int]string{
    SUCCESS:        "ok",
    ERROR:          "fail",
    INVALID_PARAMS: "Bad request params",
  }
  */
  ginerror.RegisterErrors(yourErrors.MsgFlags)
  // use gin-error-handler middleware
  r.Use(ginerror.ErrorHandle(errWriter))
  r.GET("/ping", func(c *gin.Context) {
    if "An error occurred" {
      /*
        You can implement this GenError yourself.
        but the error structure must be like this:
        type Error struct {
          StatusCode int
          Msg        string
          Code       int
        }
      */
  	  panic(ginerror.GenError(http.StatusInternalServerError,myErrors.ERROR))
    }
    c.JSON(200, gin.H{
       "message": "ok",
      })
  })
  r.Run()
}
```
