package lua4go

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"

	"time"

	"strings"

	"github.com/qxnw/lua4go/core"
	"github.com/yuin/gopher-lua"
)

var counter int32

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

func printCounter(v int32) {
	atomic.AddInt32(&counter, v)
	//if v > 0 {
	//fmt.Println("+", atomic.LoadInt32(&counter))
	//} else {
	//	fmt.Println("-", atomic.LoadInt32(&counter))
	//}
}

//NewLuaEngine 初始化lua引擎
func NewLuaEngine(script string, binder IBinder) (engine *LuaEngine, err error) {
	engine = &LuaEngine{script: script, binder: binder}
	printCounter(1)
	engine.state = lua.NewState()
	err = engine.init(script, binder)
	if err != nil {
		//	if engine.state != nil {
		//engine.state.Close()
		//}
		return
	}
	return

}
func (e *LuaEngine) init(script string, binder IBinder) (err error) {
	if err = binder.Bind(e.state); err != nil {
		return
	}
	err = e.state.DoFile(script)
	if err != nil {
		err = fmt.Errorf("脚本不存在或语法错误:%s,%+v", script, err)
		printCounter(-1)
		e.state.Close()
		return
	}
	main := e.state.GetGlobal("main")
	if main == lua.LNil {
		err = fmt.Errorf("未找到main函数:%s", script)
		return
	}
	return
}
func (e *LuaEngine) runException(context *Context, err error) {
	context.Logger.Error(err)
	printCounter(-1)
	e.state.Close()
	e.init(e.script, e.binder)
}

//Call 初始化脚本参数，并执行脚本
func (e *LuaEngine) Call(context *Context) (result []string, params map[string]string, err error) {
	if e.state == nil {
		err = fmt.Errorf("脚本不存在或语法错误:%s", e.script)
		return
	}
	defer luaRecover(context.Logger)
	startTime := time.Now()

	e.state.SetGlobal("__context__", core.New(e.state, context))
	inputData, err := json2LuaTable(e.state, context.Input, context.Logger)
	if err != nil {
		err = fmt.Errorf("脚本输入参数转换失败:%v", err)
		return
	}
	//	context.Logger.Infof("----开始执行脚本:%s", e.script)
	values, err := callMain(e.state, inputData, context.Logger)
	if err != nil {
		err = fmt.Errorf("脚本执行异常,%+v:%+v", time.Since(startTime), err)
		e.runException(context, err)
		return nil, nil, err
	}
	params = make(map[string]string)
	result = []string{}
	for _, lv := range values {
		if lv == nil || lv == lua.LNil {
			result = append(result, "")
			continue
		}
		switch lv.(type) {
		case lua.LString:
			txt := lv.String()
			if strings.HasPrefix(txt, "err:") {
				err = fmt.Errorf("脚本返回异常(%s)(%+v) %v", e.script, time.Since(startTime), lv)
				context.Logger.Error(err)
				return nil, nil, err
			}
			result = append(result, txt)
		case lua.LNumber:
			rvalue := fmt.Sprintf("%f", lv)
			if rvalue == "NaN" || rvalue == "+Inf" {
				err = fmt.Errorf("脚本返回值错误(%s)(%+v)，只支持字符串,数字和table,%v", e.script, time.Since(startTime), lv)
				e.runException(context, err)
				return
			}
			result = append(result, lv.String())
		case *lua.LTable:
			var data string
			data, params, err = luaTable2Json(lv.(*lua.LTable), context.Logger)
			if err != nil {
				err = fmt.Errorf("脚本返回结果错误(%s)(%+v),解析失败:%v", e.script, time.Since(startTime), err)
				e.runException(context, err)
				return nil, nil, err
			}
			result = append(result, data)
		default:
			err = fmt.Errorf("脚本返回值错误(%s)(%+v)，只支持字符串,数字和table:%s(current:%v)", e.script, time.Since(startTime), e.script, lv)
			e.runException(context, err)
			return
		}
	}
	params, err = getResponse(params, e.state, context.Logger)
	if err != nil {
		err = fmt.Errorf("脚本_header参数解析失败(%s)(%+v),解析失败:%v", e.script, time.Since(startTime), err)
		e.runException(context, err)
		return nil, nil, err
	}
	//	context.Logger.Infof("----脚本执行完成(%s)(%+v)", e.script, time.Since(startTime))
	return
}

//Close 关闭脚本引擎
func (e *LuaEngine) Close() {
	printCounter(-1)
	e.state.Close()
}

func callMain(ls *lua.LState, inputValue lua.LValue, log Logger) (rt []lua.LValue, er error) {
	defer func() {
		if r := recover(); r != nil {
			er = fmt.Errorf("%+v,%s", r, string(debug.Stack()))
		}
	}()
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
