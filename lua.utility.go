package lua4go

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lua4go/core"
	"github.com/yuin/gopher-lua"
)

func luaRecover(log logger.ILogger) {
	if r := recover(); r != nil {
		log.Fatal(r, string(debug.Stack()))
	}
}
func addPackages(l *lua.LState, paths ...string) (err error) {
	if paths == nil || len(paths) == 0 {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("addPackages:", r)
		}
	}()
	for _, v := range paths {
		pk := `local p = [[` + strings.Replace(v, "//", "/", -1) + `]]
local m_package_path = package.path
package.path = string.format('%s;%s/?.lua;%s/?.luac;%s/?.dll',
	m_package_path, p,p,p)`

		err = l.DoString(pk)
		if err != nil {
			return err
		}
	}

	return
}

func json2LuaTable(L *lua.LState, input string, log logger.ILogger) (inputValue lua.LValue) {
	defer luaRecover(log)
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		inputValue = lua.LString(input)
		return
	}
	tb := L.NewTable()
	for k, v := range data {
		tb.RawSetString(k, json2LuaTableValue(L, v, log))
	}
	return tb
}
func json2LuaTableValue(L *lua.LState, value interface{}, log logger.ILogger) (inputValue lua.LValue) {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Slice:
		nvalue := value.([]interface{})
		tb := L.NewTable()
		for k, v := range nvalue {
			tb.RawSetInt(k, json2LuaTableValue(L, v, log))
		}
		return tb
	case reflect.Map:
		nvalue := value.(map[string]interface{})
		tb := L.NewTable()
		for k, v := range nvalue {
			tb.RawSetString(k, json2LuaTableValue(L, v, log))
		}
		return tb
	default:
		inputValue = core.New(L, value)
	}
	return
}

func getResponse(L *lua.LState) (r map[string]string) {
	fields := map[string]string{
		"content_type":   "Content-Type",
		"charset":        "Charset",
		"original":       "_original",
		"raw":            "_original",
		"__raw__":        "_original",
		"__set_cookie__": "_cookies",
	}
	r = make(map[string]string)
	response := L.GetGlobal("response")
	if response != lua.LNil {
		for i, v := range fields {
			fied := L.GetField(response, i)
			if fied == lua.LNil {
				continue
			}
			r[v] = fied.String()
		}
	}
	for i, v := range fields {
		fied := L.GetGlobal(i)
		if fied == lua.LNil {
			continue
		}
		r[v] = fied.String()
	}
	return
}

func luaTable2Json(L *lua.LState, inputValue lua.LValue, log logger.ILogger) (json string) {
	defer luaRecover(log)
	L.Pop(L.GetTop())
	xjson := L.GetGlobal("xjson")
	if xjson.String() == "nil" {
		fmt.Println("not find xjson")
		json = inputValue.String()
		return
	}
	encode := L.GetField(xjson, "encode")
	if encode == nil {
		fmt.Println("not find xjson.encode")
		json = inputValue.String()
		return
	}
	block := lua.P{
		Fn:      encode,
		NRet:    1,
		Protect: true,
	}
	er := L.CallByParam(block, inputValue)
	if er != nil {
		fmt.Println(er)
		json = inputValue.String()
	} else {
		json = L.Get(-1).String()
	}
	L.Pop(L.GetTop())
	return
}
