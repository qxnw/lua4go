package bind

import (
	"errors"
	"fmt"

	"time"

	"github.com/qxnw/lib4go/rpc"
	"github.com/qxnw/lua4go/script"
	"github.com/yuin/gopher-lua"
)

type rpcInvoker interface {
	//Request 发送请求
	Request(service string, input map[string]string, failFast bool) (status int, result string, param map[string]string, err error)
	AsyncRequest(service string, input map[string]string, failFast bool) rpc.IRPCResponse
	WaitWithFailFast(callback func(string, int, string, error), timeout time.Duration, rs ...rpc.IRPCResponse) error
}

func moduleRPC(ls *lua.LState) (rpcInvoker, error) {
	context, err := getContext(ls)
	if err != nil {
		return nil, err
	}
	if r, ok := context.Data["__func_rpc_invoker_"]; ok {
		return r.(rpcInvoker), nil
	}
	return nil, errors.New("未找到rpc invoker")
}
func moduleRPCRequest(ls *lua.LState) int {
	service := ls.CheckString(1)
	input := ls.CheckTable(2)
	failFast := false
	if ls.GetTop() > 2 {
		failFast = ls.CheckBool(3)
	}
	rpc, err := moduleRPC(ls)
	if err != nil {
		return pushValues(ls, "", 500, err)
	}
	s, r, p, err := rpc.Request(service, getRPCMap(ls, input), failFast)
	return pushValues(ls, r, s, p, err)
}
func moduleRPCWait(ls *lua.LState) int {
	callBack := ls.CheckFunction(1)
	timeout := ls.CheckInt64(2)
	actions, err := globalRPCAction(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	if len(actions) == 0 {
		return pushValues(ls, errors.New("未指定等待的异常请求"))
	}
	rpc, err := moduleRPC(ls)
	if err != nil {
		return pushValues(ls, err)
	}
	err = rpc.WaitWithFailFast(func(svs string, s int, r string, err error) {
		co, cannel := ls.NewThread()
		defer cannel()
		fmt.Println("4.1", "svs:", svs, "s:", s, "r:", r, "err:", err, "callback:", callBack, "ls:", ls)
		ls.Resume(co, callBack, script.New(ls, svs), script.New(ls, s), script.New(ls, r), script.New(ls, err))
		fmt.Println("4.2")
	}, time.Duration(timeout)*time.Second, actions...)

	return pushValues(ls, err)
}
func globalRPCAction(ls *lua.LState) (params []rpc.IRPCResponse, err error) {
	c := ls.GetTop()
	params = make([]rpc.IRPCResponse, 0, c)
	for i := 3; i <= c; i++ {
		value := ls.CheckUserData(i)
		if lk, ok := value.Value.(rpc.IRPCResponse); ok {
			params = append(params, lk)
		} else {
			return nil, errors.New("输入参数必须为IRPCResponse类型")
		}
	}
	ls.Pop(c)
	return
}

func getRPCMap(ls *lua.LState, input *lua.LTable) map[string]string {
	data := getMapParams(input)
	cxt, err := getContext(ls)
	if err == nil {
		data["hydra_sid"] = fmt.Sprintf("%v", cxt.Data["hydra_sid"])
	}
	return data
}
