package bind

import (
	"fmt"
	"strings"
	"testing"

	"github.com/qxnw/lib4go/file"
	"github.com/qxnw/lib4go/net/http"
	"github.com/qxnw/lib4go/utility"
	"github.com/qxnw/lua4go/core"
	"github.com/yuin/gopher-lua"
)

func FuncTest(l *lua.LState) int {
	ud := l.NewUserData()
	ud.Value = http.NewHTTPClient()
	l.SetMetatable(ud, l.GetTypeMetatable("http"))
	l.Push(ud)
	return 1
}

func MethodsTest(l *lua.LState) int {
	// l.Push(core.New(l, "method"))
	ud := l.CheckUserData(1)
	v, ok := ud.Value.(*http.HTTPClient)
	if !ok {
		return 0
	}
	url := l.CheckString(2)
	params := l.CheckString(3)
	a, b, c := v.Get(url, params)
	l.Push(core.New(l, a))
	l.Push(core.New(l, b))
	l.Push(core.New(l, c))
	return 3
}

func GlobalFuncTest(l *lua.LState) int {
	l.Push(core.New(l, "global"))
	return 1
}

func ModulesTest(l *lua.LState) int {
	l.Push(core.New(l, "modules"))
	return 1
}

/*
// TestNewLuaEngine 测试构建一个lua引擎
func TestNewLuaEngine(t *testing.T) {
	packages := []string{""}
	binderTypes := []*TypeBinder{&TypeBinder{
		Name: "http",
		NewFunc: map[string]lua.LGFunction{
			"new": FuncTest,
		},
		Methods: map[string]lua.LGFunction{
			"get": MethodsTest,
		}}}
	globalFunc := map[string]lua.LGFunction{
		"global": GlobalFuncTest,
	}

	modules := map[string]map[string]lua.LGFunction{
		"modules": map[string]lua.LGFunction{
			"test": ModulesTest,
		},
	}

	// 正常加载
	filePath := file.GetAbs("./lua_test_script/test.lua")
	binder := &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	_, err := NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}

	// 没有main函数
	filePath = file.GetAbs("./lua_test_script/without_main_test.lua")
	binder = &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	_, err = NewLuaEngine(filePath, binder)
	if !strings.EqualFold(err.Error(), "未找到main函数") {
		t.Errorf("test fail : %v", err)
	}

	// 语法有误
	filePath = file.GetAbs("./lua_test_script/err_test.lua")
	_, err = NewLuaEngine(filePath, binder)
	binder = &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	if err == nil {
		t.Errorf("test fail")
	}

	// Binder没有参数
	filePath = file.GetAbs("./lua_test_script/test.lua")
	binder = &Binder{}
	_, err = NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}

	// Binder缺少Packages
	filePath = file.GetAbs("./lua_test_script/test.lua")
	binder = &Binder{Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	_, err = NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}

	// Binder缺少Types
	filePath = file.GetAbs("./lua_test_script/test.lua")
	binder = &Binder{Packages: packages, GlobalFunc: globalFunc, Modeules: modules}
	_, err = NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}

	// Binder缺少GlobalFunc
	filePath = file.GetAbs("./lua_test_script/test.lua")
	binder = &Binder{Packages: packages, Types: binderTypes, Modeules: modules}
	_, err = NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}

	// Binder缺少module
	filePath = file.GetAbs("./lua_test_script/test.lua")
	binder = &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc}
	_, err = NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}

	// 传入的文件路径不对
	binder = &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc}
	_, err = NewLuaEngine("/home/err.lua", binder)
	if err == nil {
		t.Errorf("test fail")
	}
}
*/
// TestCall 测试使用lua引擎执行lua脚本
func TestCall(t *testing.T) {
	packages := []string{"/home/champly/xlib"}
	binderTypes := []TypeBinder{
		TypeBinder{
			Name: "binder",
			NewFunc: map[string]lua.LGFunction{
				"new": FuncTest,
			},
			Methods: map[string]lua.LGFunction{
				"method": MethodsTest,
			},
		},
	}
	globalFunc := map[string]lua.LGFunction{
		"global": GlobalFuncTest,
	}

	modules := map[string]map[string]lua.LGFunction{
		"modules": map[string]lua.LGFunction{
			"test": ModulesTest,
		},
	}

	id := utility.GetGUID()
	context := &Context{Session: id[0:8], LoggerName: "luaEngine", Input: `{"id":0}`}

	// 正常加载
	filePath := file.GetAbs("./lua_test_script/test.lua")
	binder := &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	e, err := NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}

	result, _, err := e.Call(context)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if !strings.EqualFold("result", result[0]) {
		t.Errorf("test fail actual:%s\texcept:%s", result, "result")
	}

	// 没有main函数
	filePath = file.GetAbs("./lua_test_script/without_main_test.lua")
	binder = &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	e, err = NewLuaEngine(filePath, binder)
	if err == nil {
		t.Errorf("test fail")
	}

	result, _, err = e.Call(context)
	if err == nil {
		t.Errorf("test fail")
	}

	// 读取返回值
	filePath = file.GetAbs("./lua_test_script/response.lua")
	binder = &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	e, err = NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}

	// 传入string类型
	context = &Context{Session: id[0:8], LoggerName: "luaEngine", Input: `123`}
	result, params, err := e.Call(context)
	if err != nil {
		t.Errorf("test fail %v", err)
		return
	}
	if params["Charset"] != "utf-8" {
		t.Errorf("test fail actual : %s\t except : %s", params["Charset"], "utf-8")
	}
	if result[0] != "123" {
		t.Errorf("test fail actual : %s\t except : %s", result[0], "123")
	}

	// 传入json
	context = &Context{Session: id[0:8], LoggerName: "luaEngine", Input: `{"id":123}`}
	result, params, err = e.Call(context)
	if err != nil {
		t.Errorf("test fail %v", err)
		return
	}
	if params["Charset"] != "utf-8" {
		t.Errorf("test fail actual : %s\t except : %s", params["Charset"], "utf-8")
	}

	// 返回值第一个参数是302
	filePath = file.GetAbs("./lua_test_script/return_302.lua")
	binder = &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	e, err = NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	result, _, err = e.Call(context)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if len(result) != 1 {
		t.Errorf("test fail : %v", result)
	}

	// 返回值第一个参数是nil
	filePath = file.GetAbs("./lua_test_script/return_nil.lua")
	binder = &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	e, err = NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	result, params, err = e.Call(context)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if len(result) != 0 {
		t.Errorf("test fail : %v", result)
	}

	// 返回值第二个参数是nil
	filePath = file.GetAbs("./lua_test_script/return_other_nil.lua")
	binder = &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	e, err = NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	result, _, err = e.Call(context)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if len(result) != 1 {
		t.Errorf("test fail : %v", result)
	}

	// 返回值有两个
	filePath = file.GetAbs("./lua_test_script/return_two_params.lua")
	binder = &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	e, err = NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	result, params, err = e.Call(context)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if len(result) != 1 {
		t.Errorf("test fail : %v", result)
	}
}

// TestCloseEngine 测试关闭引擎
func TestCloseEngine(t *testing.T) {
	packages := []string{""}
	binderTypes := []TypeBinder{TypeBinder{Name: "http",
		NewFunc: map[string]lua.LGFunction{
			"new": FuncTest,
		},
		Methods: map[string]lua.LGFunction{
			"method": MethodsTest,
		}}}
	globalFunc := map[string]lua.LGFunction{
		"global": GlobalFuncTest,
	}

	modules := map[string]map[string]lua.LGFunction{
		"modules": map[string]lua.LGFunction{
			"test": ModulesTest,
		},
	}

	id := utility.GetGUID()
	context := &Context{Session: id[0:8], LoggerName: "luaEngine"}

	// 正常加载
	filePath := file.GetAbs("./lua_test_script/test.lua")
	binder := &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	e, err := NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}

	result, _, err := e.Call(context)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if !strings.EqualFold("result", result[0]) {
		t.Errorf("test fail actual:%s\texcept:%s", result, "result")
	}
	// 关闭lua引擎
	e.Close()

	// 测试是否关闭成功
	result, _, err = e.Call(context)
	fmt.Println(err)
	fmt.Println(result)
}

// TestMethod 测试lua脚本调用go的方法
func TestMethod(t *testing.T) {
	packages := []string{"/home/champly/xlib"}
	binderTypes := []TypeBinder{TypeBinder{
		Name: "http",
		NewFunc: map[string]lua.LGFunction{
			"new": FuncTest,
		},
		Methods: map[string]lua.LGFunction{
			"get": MethodsTest,
		}}}
	globalFunc := map[string]lua.LGFunction{
		"global": GlobalFuncTest,
	}

	modules := map[string]map[string]lua.LGFunction{
		"modules": map[string]lua.LGFunction{
			"test": ModulesTest,
		},
	}

	// 正常加载
	filePath := file.GetAbs("./lua_test_script/http.lua")
	binder := &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	e, err := NewLuaEngine(filePath, binder)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}

	id := utility.GetGUID()
	context := &Context{Session: id[0:8], LoggerName: "luaEngine", Input: `{"id":0}`}
	result, params, err := e.Call(context)
	if err != nil {
		t.Errorf("test fail : %v", err)
	}
	if !strings.EqualFold("utf-8", params["Charset"]) {
		t.Errorf("test fail actual:%s\texcept:%s", params["Charset"], "utf-8")
	}

	if !strings.EqualFold(`{"status":200,"data":"data1"}`, result[0]) {
		t.Errorf("test fail actual:%s\texcept:%s", result, `{"status":200,"data":"data1"}`)
	}

}
