package bind

import (
	"fmt"
	"strings"

	"github.com/qxnw/lib4go/file"
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
	Types      []*TypeBinder
	GlobalFunc map[string]lua.LGFunction
	Modules    map[string]map[string]lua.LGFunction
}

//NewDefault 构建默认binder
func NewDefault(pkg ...string) *Binder {
	binder := &Binder{}
	binder.Packages = getPackagePaths(pkg...)
	binder.GlobalFunc = getGlobal()
	binder.Modules = getModules()
	binder.Types = getTypes()
	return binder
}

//New 添加空的绑定函数
func New(pkg string) *Binder {
	binder := &Binder{}
	binder.Packages = getPackagePaths(pkg)
	binder.GlobalFunc = make(map[string]lua.LGFunction)
	binder.Modules = make(map[string]map[string]lua.LGFunction)
	binder.Types = []*TypeBinder{}
	return binder
}

//AppendTypes 添加自定义类型
func (b *Binder) AppendTypes(tps []*TypeBinder) {
	b.Types = append(b.Types, tps...)
}

//AppendGlobalFuncs 添加系统全局函数
func (b *Binder) AppendGlobalFuncs(fns map[string]lua.LGFunction) {
	for k, v := range fns {
		b.GlobalFunc[k] = v
	}
}

//AppendModules 添加模块
func (b *Binder) AppendModules(mds map[string]map[string]lua.LGFunction) {
	for k, v := range mds {
		b.Modules[k] = v
	}
}

//Bind 将全局函数，类型，模块等绑定到引擎
func (b *Binder) Bind(l *lua.LState) (err error) {
	err = addPackages(l, b.Packages...)
	if err != nil {
		return
	}

	for k, v := range b.Modules {
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
func getPackagePaths(p ...string) (pkgs []string) {
	pkgs = make([]string, 0, 0)
	for _, one := range p {
		ps := strings.Split(one, ";")
		for _, v := range ps {
			if v != "" {
				pkgs = append(pkgs, v)
			}
		}
	}
	return

}

func addPackages(l *lua.LState, paths ...string) (err error) {
	if paths == nil || len(paths) == 0 {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = err.(error)
		}
	}()
	for _, v := range paths {
		p, err := file.GetAbs(v)
		if err != nil {
			return fmt.Errorf("pkg path not exist :%s(err:%v)", v, err)
		}
		pk := `local p = [[` + strings.Replace(p, "//", "/", -1) + `]]
local m_package_path = package.path
package.path = string.format('%s;%s/?.lua;%s/?.luac;%s/?.dll',
	m_package_path, p,p,p)`

		err = l.DoString(pk)
		if err != nil {
			return err
		}
	}

	return
}
