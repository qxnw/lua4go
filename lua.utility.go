package lua4go

import (
	"encoding/json"
	"reflect"
	"runtime/debug"

	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lua4go/core"
	"github.com/yuin/gopher-lua"
)

func luaRecover(log logger.ILogger) {
	if r := recover(); r != nil {
		log.Fatal(r, string(debug.Stack()))
	}
}

func json2LuaTable(L *lua.LState, input string, log logger.ILogger) (inputValue lua.LValue, err error) {
	defer luaRecover(log)
	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(input), &data)
	if err != nil {
		return
	}
	tb := L.NewTable()
	for k, v := range data {
		tb.RawSetString(k, json2LuaTableValue(L, v, log))
	}
	inputValue = tb
	return
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
func getValue(L *lua.LState, obj lua.LValue, key string) (lv lua.LValue) {
	if obj == nil {
		lv = L.GetGlobal(key)
		return
	}
	lv = L.GetField(obj, key)
	if lv == lua.LNil {
		lv = L.GetGlobal(key)
		return
	}
	return lv
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
	for i, v := range fields {
		fied := getValue(L, response, i)
		if fied == lua.LNil {
			continue
		}
		r[v] = fied.String()
	}
	return
}
func luaTable2Json(tb *lua.LTable, log logger.ILogger) (s string, err error) {
	data, err := luaTable2Map(tb, log)
	if err != nil {
		return
	}
	buffer, err := json.Marshal(&data)
	if err != nil {
		return
	}
	s = string(buffer)
	return
}

func luaTable2Map(tb *lua.LTable, log logger.ILogger) (data map[string]interface{}, err error) {
	defer luaRecover(log)
	data = make(map[string]interface{})
	tb.ForEach(func(key lua.LValue, value lua.LValue) {
		if value == nil {
			return
		}
		if value.Type().String() == "table" {
			data[key.String()], err = luaTable2Map(value.(*lua.LTable), log)
			if err != nil {
				return
			}
		} else {
			data[key.String()] = value.String()
		}
	})
	return
}
