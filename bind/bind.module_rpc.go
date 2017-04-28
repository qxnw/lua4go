package bind

import (
	"errors"

	"github.com/yuin/gopher-lua"
)

type rpcInvoker interface {
	//Request 发送请求
	Request(service string, input map[string]string, failFast bool) (status int, result string, err error)
	//Query 发送请求
	Query(service string, input map[string]string, failFast bool) (status int, result string, err error)
	//Update 发送请求
	Update(service string, input map[string]string, failFast bool) (status int, err error)
	//Insert 发送请求
	Insert(service string, input map[string]string, failFast bool) (status int, err error)
	//Delete 发送请求
	Delete(service string, input map[string]string, failFast bool) (status int, err error)
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
	s, r, err := rpc.Request(service, getMapParams(input), failFast)
	return pushValues(ls, r, s, err)
}
func moduleRPCQuery(ls *lua.LState) int {
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
	s, r, err := rpc.Query(service, getMapParams(input), failFast)
	return pushValues(ls, r, s, err)
}
func moduleRPCInsert(ls *lua.LState) int {
	service := ls.CheckString(1)
	input := ls.CheckTable(2)
	failFast := false
	if ls.GetTop() > 2 {
		failFast = ls.CheckBool(3)
	}
	rpc, err := moduleRPC(ls)
	if err != nil {
		return pushValues(ls, 500, err)
	}
	s, err := rpc.Insert(service, getMapParams(input), failFast)
	return pushValues(ls, s, err)
}
func moduleRPCUpdate(ls *lua.LState) int {
	service := ls.CheckString(1)
	input := ls.CheckTable(2)
	failFast := false
	if ls.GetTop() > 2 {
		failFast = ls.CheckBool(3)
	}
	rpc, err := moduleRPC(ls)
	if err != nil {
		return pushValues(ls, 500, err)
	}
	s, err := rpc.Update(service, getMapParams(input), failFast)
	return pushValues(ls, s, err)
}
func moduleRPCDelete(ls *lua.LState) int {
	service := ls.CheckString(1)
	input := ls.CheckTable(2)
	failFast := false
	if ls.GetTop() > 2 {
		failFast = ls.CheckBool(3)
	}
	rpc, err := moduleRPC(ls)
	if err != nil {
		return pushValues(ls, 500, err)
	}
	s, err := rpc.Delete(service, getMapParams(input), failFast)
	return pushValues(ls, s, err)
}
