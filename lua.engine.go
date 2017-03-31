package lua4go

import (
	"fmt"

	"time"

	"github.com/qxnw/lua4go/core"
	"github.com/yuin/gopher-lua"
)

//IBinder 基础库绑定
type IBinder interface {
	Bind(*lua.LState) error
}

//LuaEngine 脚本引擎
type LuaEngine struct {
	binder IBinder
	script string
	state  *lua.LState
}

//NewLuaEngine 初始化lua引擎
func NewLuaEngine(script string, binder IBinder) (engine *LuaEngine, err error) {
	engine = &LuaEngine{script: script, binder: binder}
	engine.state = lua.NewState()
	if err = binder.Bind(engine.state); err != nil {
		return
	}
	err = engine.state.DoFile(script)
	if err != nil {
		err = fmt.Errorf("脚本语法错误:%s,%+v", script, err)
		engine.state.Close()
		return
	}
	main := engine.state.GetGlobal("main")
	if main == lua.LNil {
		err = fmt.Errorf("未找到main函数:%s", script)
		return
	}
	return
}

//Call 初始化脚本参数，并执行脚本
func (e *LuaEngine) Call(context *Context) (result []string, params map[string]string, err error) {
	defer luaRecover(context.Logger)
	startTime := time.Now()

	e.state.SetGlobal("__context__", core.New(e.state, context))
	inputData, err := json2LuaTable(e.state, context.Input, context.Logger)
	if err != nil {
		err = fmt.Errorf("脚本输入参数转换失败:%v", err)
		return
	}
	context.Logger.Infof("----开始执行脚本:%s", e.script)
	values, err := callMain(e.state, inputData, context.Logger)
	if err != nil {
		err = fmt.Errorf("脚本执行异常:%+v", err)
		return
	}
	result = []string{}
	for _, lv := range values {
		switch lv.(type) {
		case lua.LString:
			result = append(result, lv.String())
		case lua.LNumber:
			result = append(result, lv.String())
		case *lua.LTable:
			data, err := luaTable2Json(lv.(*lua.LTable), context.Logger)
			if err != nil {
				err = fmt.Errorf("脚本返回结果解析失败:%v", err)
				return nil, nil, err
			}
			result = append(result, data)
		default:
			err = fmt.Errorf("脚本返回值错误，只支持字符串和table:%s", e.script)
		}
	}
	params = getResponse(e.state)
	context.Logger.Infof("----完成执行脚本:%s,%+v", e.script, time.Since(startTime))
	return
}

//Close 关闭脚本引擎
func (e *LuaEngine) Close() {
	e.state.Close()
}

func callMain(ls *lua.LState, inputValue lua.LValue, log Logger) (rt []lua.LValue, er error) {
	defer luaRecover(log)
	ls.Pop(ls.GetTop())
	er = callMainFunc(ls, inputValue)
	if er != nil {
		return
	}
	defer ls.Pop(ls.GetTop())
	count := ls.GetTop()
	rt = make([]lua.LValue, 0, count)
	if count == 0 {
		return
	}
	value1 := ls.Get(1)
	rt = append(rt, value1)
	return
}
func callMainFunc(ls *lua.LState, args ...lua.LValue) (err error) {
	block := lua.P{
		Fn:      ls.GetGlobal("main"),
		NRet:    1,
		Protect: true,
	}
	return ls.CallByParam(block, args...)
}
