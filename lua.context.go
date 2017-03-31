package lua4go

import (
	"net/http"

	"os"

	"github.com/qxnw/lib4go/logger"
)

type HttpContext struct {
	Response http.ResponseWriter
	Request  *http.Request
}

//Context 脚本执行上下文
type Context struct {
	Logger      Logger
	Input       string
	Response    http.ResponseWriter
	HttpContext *HttpContext
	Data        map[string]string
}

func NewContext(input string) *Context {
	return NewContextWithLogger(input, NewLogger(os.Stdout))
}

//NewContextWithLogger 初始化Context
func NewContextWithLogger(input string, logger Logger) *Context {
	return &Context{Logger: logger, Input: input}
}

//NewContextHTTP  初始化Context
func NewContextHTTP(logger logger.ILogger, input string, w http.ResponseWriter,
	r *http.Request) *Context {
	return &Context{Logger: logger, Input: input, HttpContext: &HttpContext{Response: w, Request: r}}
}
