package bind

import (
	"errors"

	"github.com/arsgo/lib4go/script"
	"github.com/qxnw/lua4go"
	"github.com/yuin/gopher-lua"
)

func pushValues(ls *lua.LState, values ...interface{}) int {
	for _, v := range values {
		ls.Push(script.New(ls, v))
	}
	return len(values)
}
func globalGetLogger(ls *lua.LState) (lg lua4go.Logger, err error) {
	context, err := getContext(ls)
	if err != nil {
		return nil, err
	}
	return context.Logger, nil
}
func moduleHTTPContext(ls *lua.LState) (*lua4go.HttpContext, error) {
	context, err := getContext(ls)
	if err != nil {
		return nil, err
	}
	return context.HttpContext, nil
}
func getContext(ls *lua.LState) (*lua4go.Context, error) {
	context := ls.GetGlobal("__context__")
	if context == nil {
		return nil, errors.New("未找到脚本请求上下文变量")
	}
	data := context.(*lua.LUserData)
	hr := data.Value.(*lua4go.Context)
	return hr, nil
}
func getStringParams(ls *lua.LState, start int) (params []string) {
	c := ls.GetTop()
	params = make([]string, 0, c)
	for i := start; i <= c; i++ {
		t := ls.Get(i).Type().String()
		if t == "userdata" {
			ls.RaiseError("invalid string of function arguments (string expected, got userdata)")
		} else {
			params = append(params, ls.Get(i).String())
		}
	}
	return
}
func getMapParams(tb *lua.LTable) map[string]string {
	data := make(map[string]string)
	tb.ForEach(func(key lua.LValue, value lua.LValue) {
		if value != nil {
			data[key.String()] = value.String()
		}
	})
	return data
}
