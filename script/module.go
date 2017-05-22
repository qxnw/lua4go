package script

import "github.com/yuin/gopher-lua"

type luaModule struct {
	exports map[string]lua.LGFunction
}

func NewLuaModule(exports map[string]lua.LGFunction) luaModule {
	return luaModule{exports: exports}
}

func (l luaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), l.exports)
	L.Push(mod)
	return 1
}
