package ginerror

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"runtime"

	ginErrors "github.com/garfield-yin/gin-error-handler/errors"
	"github.com/gin-gonic/gin"
)

// Error Struct
type Error struct {
	StatusCode int
	Msg        string
	Code       int
}

func (err *Error) Error() string {
	return fmt.Sprintf("status_code:%d, msg:%s", err.StatusCode, err.Msg)
}

// RegisterErrors : register your error messages
func RegisterErrors(msgFlags map[int]string) {
	for key, value := range msgFlags {
		ginErrors.MsgFlags[key] = value
	}
}

// GenError error build, You can implement this function yourself.
func GenError(httpCode int, errCode int, msg ...string) Error {
	err := Error{
		StatusCode: httpCode,
		Code:       errCode,
	}
	if len(msg) > 0 { // your message stri
		err.Msg = msg[0]
	} else {
		err.Msg = ginErrors.GetMsg(errCode)
	}
	return err
}

func abortWithError(c *gin.Context, err Error) {
	c.JSON(err.StatusCode, gin.H{
		"code":    err.Code,
		"message": err.Msg,
	})
	c.Abort()
}

// ErrorHandle 统一捕获错误
func ErrorHandle(out io.Writer) gin.HandlerFunc {
	logger := log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(Error); ok {
					//自定义错误，业务逻辑故意抛出的，返回统一格式数据
					abortWithError(ctx, e)
					return
				}

				// error stack
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]
				httprequest, _ := httputil.DumpRequest(ctx.Request, false)
				logger.Printf("[Recovery] panic recovered:\n%s\n%s\n%s", string(httprequest), err, stack)
				// default error
				internalServerError := GenError(http.StatusInternalServerError, ginErrors.ERROR)
				abortWithError(ctx, internalServerError)
			}
		}()
		ctx.Next()
	}
}
