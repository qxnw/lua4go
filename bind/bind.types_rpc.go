package bind

import (
	"time"

	"github.com/qxnw/lib4go/rpc"
	"github.com/yuin/gopher-lua"
)

func getAsyncRpcResponseTypeBinder() *TypeBinder {
	return &TypeBinder{
		Name:    "async_rpc_response",
		NewFunc: map[string]lua.LGFunction{},
		Methods: map[string]lua.LGFunction{
			"wait": typeRPCReponseWait,
		},
	}
}
func checkRPCAsyncResponseType(L *lua.LState) rpc.IRPCResponse {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(rpc.IRPCResponse); ok {
		return v
	}
	L.RaiseError("bad argument  (RPC Invoker expected, got %s)", ud.Type().String())
	return nil
}

func typeRPCReponseWait(ls *lua.LState) int {
	response := checkRPCAsyncResponseType(ls)
	timeout := ls.CheckInt64(2)
	status, result, err := response.Wait(time.Duration(timeout) * time.Millisecond)
	if err != nil {
		return pushValues(ls, "", 500, err)
	}
	return pushValues(ls, status, result)
}
func moduleRPCAsyncDelete(ls *lua.LState) int {
	rpc, err := moduleRPC(ls)
	if err != nil {
		return pushValues(ls, "", 500, err)
	}
	service := ls.CheckString(1)
	input := ls.CheckTable(2)
	failFast := true
	if ls.GetTop() > 2 {
		failFast = ls.CheckBool(3)
	}
	d := rpc.AsyncDelete(service, getRPCMap(ls, input), failFast)
	ud := ls.NewUserData()
	ud.Value = d
	ls.SetMetatable(ud, ls.GetTypeMetatable("async_rpc_response"))
	ls.Push(ud)
	return 1
}
func moduleRPCAsyncInsert(ls *lua.LState) int {
	rpc, err := moduleRPC(ls)
	if err != nil {
		return pushValues(ls, "", 500, err)
	}
	service := ls.CheckString(1)
	input := ls.CheckTable(2)
	failFast := true
	if ls.GetTop() > 2 {
		failFast = ls.CheckBool(3)
	}
	d := rpc.AsyncInsert(service, getRPCMap(ls, input), failFast)
	ud := ls.NewUserData()
	ud.Value = d
	ls.SetMetatable(ud, ls.GetTypeMetatable("async_rpc_response"))
	ls.Push(ud)
	return 1
}
func moduleRPCAsyncQuery(ls *lua.LState) int {
	rpc, err := moduleRPC(ls)
	if err != nil {
		return pushValues(ls, "", 500, err)
	}
	service := ls.CheckString(1)
	input := ls.CheckTable(2)
	failFast := true
	if ls.GetTop() > 2 {
		failFast = ls.CheckBool(3)
	}
	d := rpc.AsyncQuery(service, getRPCMap(ls, input), failFast)
	ud := ls.NewUserData()
	ud.Value = d
	ls.SetMetatable(ud, ls.GetTypeMetatable("async_rpc_response"))
	ls.Push(ud)
	return 1
}

func moduleRPCAsyncRequest(ls *lua.LState) int {
	rpc, err := moduleRPC(ls)
	if err != nil {
		return pushValues(ls, "", 500, err)
	}
	service := ls.CheckString(1)
	input := ls.CheckTable(2)
	failFast := true
	if ls.GetTop() > 2 {
		failFast = ls.CheckBool(3)
	}
	d := rpc.AsyncRequest(service, getRPCMap(ls, input), failFast)
	ud := ls.NewUserData()
	ud.Value = d
	ls.SetMetatable(ud, ls.GetTypeMetatable("async_rpc_response"))
	ls.Push(ud)
	return 1
}
func moduleRPCAsyncUpdate(ls *lua.LState) int {
	rpc, err := moduleRPC(ls)
	if err != nil {
		return pushValues(ls, "", 500, err)
	}
	service := ls.CheckString(1)
	input := ls.CheckTable(2)
	failFast := true
	if ls.GetTop() > 2 {
		failFast = ls.CheckBool(3)
	}
	d := rpc.AsyncUpdate(service, getRPCMap(ls, input), failFast)
	ud := ls.NewUserData()
	ud.Value = d
	ls.SetMetatable(ud, ls.GetTypeMetatable("async_rpc_response"))
	ls.Push(ud)
	return 1
}
