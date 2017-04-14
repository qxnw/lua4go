package bind

import lua "github.com/yuin/gopher-lua"

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
