package bind

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func moduleHTTPContextGetCookie(ls *lua.LState) int {
	key := ls.CheckString(1)
	request, err := moduleHTTPRequest(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	ck, err := request.Cookie(key)
	if err != nil {
		return pushValues(ls, err)
	}
	return pushValues(ls, ck.Value)
}

func moduleHTTPContextSetCookie(ls *lua.LState) int {
	cookies := ls.CheckString(1)
	response, err := moduleHTTPResponseWriter(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	response.Header().Add("Set-Cookie", cookies)
	return pushValues(ls)
}
func moduleHTTPContextSetContentType(ls *lua.LState) int {
	value := ls.CheckString(1)
	response, err := moduleHTTPResponseWriter(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	response.Header().Set("Content-Type", value)
	return pushValues(ls)
}

//__func_var_get_

func getFuncVarGet(ls *lua.LState, t, n string) (v string, err error) {
	context, err := getContext(ls)
	if err != nil {
		return
	}
	if r, ok := context.Data["__func_var_get_"]; ok {
		if fun, ok := r.(func(string, string) (string, error)); ok {
			v, err = fun(t, n)
			return
		}
	}
	return "", fmt.Errorf("未找到__func_var_get_")
}

func moduleGetVar(ls *lua.LState) int {
	value1 := ls.CheckString(1)
	value2 := ls.CheckString(2)
	value, err := getFuncVarGet(ls, value1, value2)
	if err != nil {
		return pushValues(ls, "", err)
	}
	return pushValues(ls, value)
}

func moduleGetBody(ls *lua.LState) int {
	context, err := getContext(ls)
	if err != nil {
		return pushValues(ls, "", err)
	}
	value := ls.CheckString(1)
	if r, ok := context.Data["__func_body_get_"]; ok {
		if fun, ok := r.(func(string) (string, error)); ok {
			if value, err := fun(value); err == nil {
				return pushValues(ls, value)
			}
			return pushValues(ls, "", err)
		}
	}
	return pushValues(ls, "", "不支持")
}
func moduleContexSetCharset(ls *lua.LState) int {
	value := ls.CheckString(1)
	response, err := moduleHTTPResponseWriter(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	response.Header().Set("Charset", value)
	return pushValues(ls)
}

func moduleContexSetHeader(ls *lua.LState) int {
	key := ls.CheckString(1)
	value := ls.CheckString(2)
	response, err := moduleHTTPResponseWriter(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	response.Header().Set(key, value)
	return pushValues(ls)
}
