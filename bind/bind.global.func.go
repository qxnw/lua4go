package bind

import (
	"time"

	"fmt"

	"github.com/qxnw/lib4go/utility"
	"github.com/yuin/gopher-lua"
)

//用于获取当前运行时的输入参数
func globalGetParams(ls *lua.LState) (params []interface{}) {
	c := ls.GetTop()
	params = make([]interface{}, 0, c)
	for i := 1; i <= c; i++ {
		value := ls.Get(i)
		switch value.(type) {
		case *lua.LUserData:
			params = append(params, ls.CheckUserData(i).Value)
		default:
			params = append(params, value.String())
		}
	}
	ls.Pop(c)
	return
}

//获取guid
func globalGUID(ls *lua.LState) int {
	return pushValues(ls, utility.GetGUID())
}

//获取日志组件的info函数
func globalInfo(ls *lua.LState) int {
	params := globalGetParams(ls)
	if len(params) == 0 {
		return pushValues(ls)
	}
	lg, err := globalGetLogger(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	lg.Info(params...)
	return pushValues(ls)
}

//获取日志组件的infof函数
func globalInfof(ls *lua.LState) int {
	params := globalGetParams(ls)
	if len(params) <= 1 {
		return pushValues(ls, fmt.Errorf(`bad argument #%v to %v (%v expected, got %v)`, "1", "infof", "string", "nil"))
	}
	lg, err := globalGetLogger(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	lg.Infof(params[0].(string), params[1:]...)
	return pushValues(ls)
}

func globalError(ls *lua.LState) int {
	params := globalGetParams(ls)
	if len(params) == 0 {
		return pushValues(ls)
	}
	lg, err := globalGetLogger(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	lg.Error(params...)
	return pushValues(ls)
}
func globalErrorf(ls *lua.LState) int {
	params := globalGetParams(ls)
	if len(params) <= 1 {
		return pushValues(ls, fmt.Errorf(`bad argument #%v to %v (%v expected, got %v)`, "1", "errorf", "string", "nil"))
	}
	lg, err := globalGetLogger(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	lg.Errorf(params[0].(string), params[1:]...)
	return pushValues(ls)
}
func globalSleep(ls *lua.LState) int {
	second := ls.CheckInt(1)
	time.Sleep(time.Second * time.Duration(second))
	return pushValues(ls)
}
