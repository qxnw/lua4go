package lua4go

import (
	"encoding/json"
	"reflect"
	"runtime/debug"

	"fmt"

	"github.com/qxnw/lua4go/core"
	"github.com/yuin/gopher-lua"
)

func luaRecover(log Logger) {
	if r := recover(); r != nil {
		log.Error(r, string(debug.Stack()))
	}
}

func json2LuaTable(L *lua.LState, input string, log Logger) (inputValue lua.LValue, err error) {
	defer luaRecover(log)
	if input == "" {
		input = "{}"
	}
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
func json2LuaTableValue(L *lua.LState, value interface{}, log Logger) (inputValue lua.LValue) {
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
	if obj == nil || obj == lua.LNil {
		lv = L.GetGlobal(key)
		return
	}
	lv = L.GetField(obj, key)
	if lv == lua.LNil || lv == nil {
		lv = L.GetGlobal(key)
		return
	}
	return lv
}

func getResponse(params map[string]string, L *lua.LState, log Logger) (r map[string]string, err error) {
	response := L.GetGlobal("header")
	if response == nil || response == lua.LNil {
		response = L.GetGlobal("response")
		if response != nil && response != lua.LNil {
			if h, ok := response.(*lua.LTable); ok {
				response = h.RawGet(lua.LString("header"))
			}
		}
	}
	if response == nil || response == lua.LNil {
		return params, nil
	}

	//检查header是否是luatable
	if _, ok := response.(*lua.LTable); !ok {
		return params, nil
	}

	//转换头信息
	m, err := luaTable2Map(response.(*lua.LTable), log)
	if err != nil {
		return
	}

	//填充返回参数
	for i, v := range m {
		if _, ok := params[i]; !ok && v != nil {
			params[i] = fmt.Sprintf("%v", v)
		}
	}
	return params, nil
}
func luaTable2Json(tb *lua.LTable, log Logger) (s string, m map[string]string, err error) {
	m = make(map[string]string)
	data, err := luaTable2Map(tb, log)
	if err != nil {
		return
	}
	if v, ok := data["header"]; ok {
		header, ok := v.(map[string]interface{})
		if ok {
			for k, value := range header {
				m[k] = fmt.Sprintf("%v", value)
			}
		}
		delete(data, "header")
	}
	buffer, err := json.Marshal(&data)
	if err != nil {
		return
	}
	s = string(buffer)
	return
}

func luaTable2Map(tb *lua.LTable, log Logger) (data map[string]interface{}, err error) {
	defer luaRecover(log)
	data = make(map[string]interface{})
	tb.ForEach(func(key lua.LValue, value lua.LValue) {
		if value == nil || value == lua.LNil {
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
