package lua4go

import (
	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lib4go/utility"
)

//Context 脚本执行上下文
type Context struct {
	Logger *logger.Logger
	Input  string
	Data   map[string]interface{}
}

func NewContext(input string) *Context {
	return NewContextWithLogger(input, map[string]interface{}{}, logger.GetSession("script", utility.GetGUID()))
}

//NewContextWithLogger 初始化Context
func NewContextWithLogger(input string, data map[string]interface{}, logger *logger.Logger) *Context {
	return &Context{Logger: logger, Data: data, Input: input}
}
