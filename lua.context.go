package lua4go

import (
	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lib4go/server/http"
)

//Context 脚本执行上下文
type Context struct {
	Logger      logger.ILogger
	Input       string
	HttpContext *http.Context
}

//NewContext 初始化Context
func NewContext(logger logger.ILogger, input string) *Context {
	return &Context{Logger: logger, Input: input}
}

//NewContextHTTP  初始化Context
func NewContextHTTP(logger logger.ILogger, input string, httpContext *http.Context) *Context {
	return &Context{Logger: logger, Input: input, HttpContext: httpContext}
}
