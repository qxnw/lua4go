package lua4go

import (
	"github.com/qxnw/lua4go/core"
	"github.com/yuin/gopher-lua"
)

//TypeBinder 类型绑定信息
type TypeBinder struct {
	Name    string
	NewFunc map[string]lua.LGFunction
	Methods map[string]lua.LGFunction
}

//Binder 绑定对象
type Binder struct {
	Packages   []string
	Types      []TypeBinder
	GlobalFunc map[string]lua.LGFunction
	Modeules   map[string]map[string]lua.LGFunction
}

//Bind 将全局函数，类型，模块等绑定到引擎
func (b *Binder) Bind(l *lua.LState) (err error) {
	err = addPackages(l, b.Packages...)
	if err != nil {
		return
	}

	for k, v := range b.Modeules {
		l.PreloadModule(k, core.NewLuaModule(v).Loader)
	}

	for _, v := range b.Types {
		mt := l.NewTypeMetatable(v.Name)
		l.SetGlobal(v.Name, mt)
		for i, ff := range v.NewFunc {
			l.SetField(mt, i, l.NewFunction(ff))
		}
		l.SetField(mt, "__index", l.SetFuncs(l.NewTable(), v.Methods))
	}

	for i, v := range b.GlobalFunc {
		l.SetGlobal(i, l.NewFunction(v))
	}
	return

}
