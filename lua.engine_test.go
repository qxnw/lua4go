package lua4go

import (
	"fmt"
	"strings"
	"testing"

	"github.com/qxnw/lib4go/file"
	"github.com/qxnw/lib4go/utility"
	"github.com/yuin/gopher-lua"
)

func FuncTest(l *lua.LState) int {
	return 0
}

func MethodsTest(l *lua.LState) int {
	return 1
}

func GlobalFuncTest(l *lua.LState) int {
	return 2
}

func ModulesTest(l *lua.LState) int {
	return 3
}

// TestNewLuaEngine 测试构建一个lua引擎
func TestNewLuaEngine(t *testing.T) {
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
}

// TestCall 测试使用lua引擎执行lua脚本
func TestCall(t *testing.T) {
	packages := []string{""}
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

	// 没有main函数
	filePath = file.GetAbs("./lua_test_script/without_main_test.lua")
	binder = &Binder{Packages: packages, Types: binderTypes, GlobalFunc: globalFunc, Modeules: modules}
	e, err = NewLuaEngine(filePath, binder)
	// if err != nil {
	// 	t.Errorf("test fail : %v", err)
	// }

	result, _, err = e.Call(context)
	// fmt.Println(err)
	if err == nil {
		t.Errorf("test fail")
	}
}

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
