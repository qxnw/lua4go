package script

import (
	"reflect"

	"github.com/yuin/gopher-lua"
)

func checkStruct(L *lua.LState, idx int) (ref reflect.Value, mt *Metatable, isPtr bool) {
	ud := L.CheckUserData(idx)
	ref = reflect.ValueOf(ud.Value)
	if ref.Kind() != reflect.Struct {
		if ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
			L.ArgError(idx, "expecting struct or struct pointer")
		}
		isPtr = true
	}
	mt = &Metatable{LTable: ud.Metatable.(*lua.LTable)}
	return
}

func structIndex(L *lua.LState) int {
	ref, mt, isPtr := checkStruct(L, 1)
	key := L.CheckString(2)

	if isPtr {
		if fn := mt.ptrMethod(key); fn != nil {
			L.Push(fn)
			return 1
		}
	}

	if fn := mt.method(key); fn != nil {
		L.Push(fn)
		return 1
	}

	ref = reflect.Indirect(ref)
	index := mt.fieldIndex(key)
	if index == nil {
		return 0
	}
	field := ref.FieldByIndex(index)
	if !field.CanInterface() {
		L.RaiseError("cannot interface field " + key)
	}
	switch field.Kind() {
	case reflect.Array, reflect.Struct:
		if field.CanAddr() {
			field = field.Addr()
		}
	}
	L.Push(New(L, field.Interface()))
	return 1
}

func structNewIndex(L *lua.LState) int {
	ref, mt, isPtr := checkStruct(L, 1)
	if isPtr {
		ref = ref.Elem()
	}
	key := L.CheckString(2)
	value := L.CheckAny(3)

	index := mt.fieldIndex(key)
	if index == nil {
		L.RaiseError("unknown field " + key)
	}
	field := ref.FieldByIndex(index)
	if !field.CanSet() {
		L.RaiseError("cannot set field " + key)
	}
	field.Set(lValueToReflect(value, field.Type()))
	return 0
}
