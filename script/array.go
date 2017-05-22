package script

import (
	"reflect"

	"github.com/yuin/gopher-lua"
)

func checkArray(L *lua.LState, idx int) (ref reflect.Value, mt *Metatable, isPtr bool) {
	ud := L.CheckUserData(idx)
	ref = reflect.ValueOf(ud.Value)
	if ref.Kind() != reflect.Array {
		if ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Array {
			L.ArgError(idx, "expecting array or array pointer")
		}
		isPtr = true
	}
	mt = &Metatable{LTable: ud.Metatable.(*lua.LTable)}
	return
}

func arrayIndex(L *lua.LState) int {
	ref, mt, isPtr := checkArray(L, 1)
	ref = reflect.Indirect(ref)
	key := L.CheckAny(2)

	switch converted := key.(type) {
	case lua.LNumber:
		index := int(converted)
		if index < 1 || index > ref.Len() {
			L.ArgError(2, "index out of range")
		}
		L.Push(New(L, ref.Index(index-1).Interface()))
	case lua.LString:
		if isPtr {
			if fn := mt.ptrMethod(string(converted)); fn != nil {
				L.Push(fn)
				return 1
			}
		}
		if fn := mt.method(string(converted)); fn != nil {
			L.Push(fn)
			return 1
		}
		return 0
	default:
		L.ArgError(2, "must be a number or string")
	}
	return 1
}

func arrayNewIndex(L *lua.LState) int {
	ref, _, isPtr := checkArray(L, 1)

	if !isPtr {
		L.RaiseError("invalid operation on array")
	}

	ref = ref.Elem()

	index := L.CheckInt(2)
	value := L.CheckAny(3)
	if index < 1 || index > ref.Len() {
		L.ArgError(2, "index out of range")
	}
	ref.Index(index - 1).Set(lValueToReflect(value, ref.Type().Elem()))
	return 0
}

func arrayLen(L *lua.LState) int {
	ref, _, _ := checkArray(L, 1)
	ref = reflect.Indirect(ref)
	L.Push(lua.LNumber(ref.Len()))
	return 1
}
