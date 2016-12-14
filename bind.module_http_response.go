package lua4go

import (
	lua "github.com/yuin/gopher-lua"
)

func moduleHTTPContextGetCookie(ls *lua.LState) int {
	key := ls.CheckString(1)
	context, err := moduleHTTPContext(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	ck, err := context.Request.Cookie(key)
	if err != nil {
		return pushValues(ls, err)
	}
	return pushValues(ls, ck.Value)
}

func moduleHTTPContextSetCookie(ls *lua.LState) int {
	cookies := ls.CheckString(1)
	context, err := moduleHTTPContext(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	context.Response.Header().Add("Set-Cookie", cookies)
	return pushValues(ls)
}
func moduleHTTPContextSetContentType(ls *lua.LState) int {
	value := ls.CheckString(1)
	context, err := moduleHTTPContext(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	context.Response.Header().Set("Content-Type", value)
	return pushValues(ls)
}
func moduleContexSetCharset(ls *lua.LState) int {
	value := ls.CheckString(1)
	context, err := moduleHTTPContext(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	context.Response.Header().Set("Charset", value)
	return pushValues(ls)
}
func moduleContexSetHeader(ls *lua.LState) int {
	key := ls.CheckString(1)
	value := ls.CheckString(2)
	context, err := moduleHTTPContext(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	context.Response.Header().Set(key, value)
	return pushValues(ls)
}
