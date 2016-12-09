package lua4go

import (
	"errors"
	"strings"

	"github.com/qxnw/lib4go/logger"
	"github.com/qxnw/lua4go/core"
	"github.com/yuin/gopher-lua"
)

//LuaEngine 脚本引擎
type LuaEngine struct {
	Binder *Binder
	script string
	state  *lua.LState
}

//NewLuaEngine 初始化lua引擎
func NewLuaEngine(script string, binder *Binder) (engine *LuaEngine, err error) {
	engine = &LuaEngine{script: script, Binder: binder}
	engine.state = lua.NewState()
	if err = binder.Bind(engine.state); err != nil {
		return
	}
	err = engine.state.DoFile(script)
	if err != nil {
		engine.state.Close()
		return
	}
	main := engine.state.GetGlobal("main")
	if main == lua.LNil {
		err = errors.New("未找到main函数")
		return
	}
	return
}

//Call 初始化脚本参数，并执行脚本
func (e *LuaEngine) Call(context *Context) (result []string, params map[string]string, err error) {
	log := logger.GetSession(context.LoggerName, context.Session)
	defer luaRecover(log)
	log.Info("----开始执行脚本:%s", e.script)

	e.state.SetGlobal("__session__", lua.LString(context.Session))
	e.state.SetGlobal("__logger_name__", lua.LString(context.LoggerName))
	e.state.SetGlobal("__http_context__", core.New(e.state, context.HTTPContext))

	inputData := json2LuaTable(e.state, context.Input, log)
	values, err := callMain(e.state, inputData, lua.LString(context.Body), log)
	if err != nil {
		return
	}
	result = []string{}
	for _, lv := range values {
		if strings.EqualFold(lv.Type().String(), "table") {
			result = append(result, luaTable2Json(e.state, lv, log))
		} else {
			result = append(result, lv.String())
		}
	}
	params = getResponse(e.state)
	log.Info("----完成执行脚本:%s", e.script)
	return
}

//Close 关闭脚本引擎
func (e *LuaEngine) Close() {
	e.state.Close()
}

func callMain(ls *lua.LState, inputValue lua.LValue, others lua.LValue, log logger.ILogger) (rt []lua.LValue, er error) {
	defer luaRecover(log)
	ls.Pop(ls.GetTop())
	er = callMainFunc(ls, inputValue, others)
	if er != nil {
		return
	}
	defer ls.Pop(ls.GetTop())
	rt = make([]lua.LValue, 0, ls.GetTop())
	value1 := ls.Get(1)
	if value1.String() == "nil" {
		return
	}
	rt = append(rt, value1)
	if value1.String() != "302" {
		return
	}
	value2 := ls.Get(2)
	if value2.String() == "nil" {
		return
	}
	rt = append(rt, value2)
	return
}
func callMainFunc(ls *lua.LState, args ...lua.LValue) (err error) {
	block := lua.P{
		Fn:      ls.GetGlobal("main"),
		NRet:    2,
		Protect: true,
	}
	err = ls.CallByParam(block, args...)
	return err
}
