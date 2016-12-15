package bind

import (
	"fmt"
	"time"

	"errors"

	"github.com/qxnw/lib4go/utility"
	"github.com/yuin/gopher-lua"
)

func globalGetParams(ls *lua.LState) (params []interface{}) {
	c := ls.GetTop()
	params = make([]interface{}, 0, c)
	for i := 1; i <= c; i++ {
		t := ls.Get(i).Type().String()
		if t == "userdata" {
			params = append(params, fmt.Sprintf("%+v", ls.CheckUserData(i).Value))
		} else {
			params = append(params, ls.Get(i).String())
		}
	}
	return
}

func globalGUID(ls *lua.LState) int {
	return pushValues(ls, utility.GetGUID())
}

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
func globalInfof(ls *lua.LState) int {
	params := globalGetParams(ls)
	if len(params) <= 1 {
		return pushValues(ls, errors.New("输入参数个数有误"))
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
		return pushValues(ls, errors.New("输入参数个数有误"))
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
