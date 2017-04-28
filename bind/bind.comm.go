package bind

import (
	"errors"
	"net/http"

	"github.com/arsgo/lib4go/script"
	"github.com/qxnw/lua4go"
	"github.com/yuin/gopher-lua"
)

func pushValues(ls *lua.LState, values ...interface{}) int {
	for _, v := range values {
		if err, ok := v.(error); ok {
			ls.Push(script.New(ls, "err:"+err.Error()))
			continue
		}
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
func moduleHTTPResponseWriter(ls *lua.LState) (http.ResponseWriter, error) {
	context, err := getContext(ls)
	if err != nil {
		return nil, err
	}
	if r, ok := context.Data["__func_http_response_"]; ok {
		return r.(http.ResponseWriter), nil
	}
	return nil, errors.New("未找到http.ResponseWriter")
}

func moduleHTTPRequest(ls *lua.LState) (*http.Request, error) {
	context, err := getContext(ls)
	if err != nil {
		return nil, err
	}
	if r, ok := context.Data["__func_http_request_"]; ok {
		return r.(*http.Request), nil
	}
	return nil, errors.New("未找到http.ResponseWriter")
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
func getIMapParams(tb *lua.LTable) map[string]interface{} {
	data := make(map[string]interface{})
	tb.ForEach(func(key lua.LValue, value lua.LValue) {
		if value != nil && value != lua.LNil {
			data[key.String()] = value.String()
		}
	})
	return data
}
func toLuaTable(ls *lua.LState, input []interface{}) (tb *lua.LTable) {
	tb = ls.NewTable()
	for i, v := range input {
		tb.RawSetInt(i+1, script.New(ls, v))
	}
	return
}
