package lua4go

import (
	"reflect"
	"testing"

	"github.com/yuin/gopher-lua"
)

type Binder struct {
}

func (b *Binder) Bind(*lua.LState) error {
	return nil
}

type watcher struct {
	callback func()
}

func (w *watcher) Append(f string) error {
	return nil
}

func TestVM1(t *testing.T) {
	vm := NewLuaVM(&Binder{}, 1, 1)
	wc := &watcher{callback: vm.Reload}
	vm.watcher = wc
	//检查初始值
	expect(t, vm.version, int32(100))
	expect(t, vm.cache.Count(), 1)
	expect(t, vm.scripts.Count(), 0)

	//检查已加载的脚本数
	err := vm.PreLoad("./test/t1.lua")
	expect(t, err, nil)
	expect(t, vm.scripts.Count(), 1)
	expect(t, vm.version, int32(100))
	//重复加载脚本
	err = vm.PreLoad("./test/t1.lua")
	expect(t, err, nil)
	expect(t, vm.scripts.Count(), 1)
	expect(t, vm.version, int32(100))

	//加载不存在脚本
	err = vm.PreLoad("./test/t1_not_exist.lua")
	refute(t, err, nil)
	expect(t, vm.scripts.Count(), 1)
	expect(t, vm.version, int32(100))

}

func TestVM2(t *testing.T) {
	vm := NewLuaVM(&Binder{}, 1, 1)
	wc := &watcher{callback: vm.Reload}
	vm.watcher = wc
	//检查初始值
	expect(t, vm.version, int32(100))
	expect(t, vm.cache.Count(), 1)
	expect(t, vm.scripts.Count(), 0)

	//检查已加载的脚本数
	err := vm.PreLoad("./test/t1.lua")
	expect(t, err, nil)
	expect(t, vm.scripts.Count(), 1)

	//回调后检查缓存引擎及脚本数
	for i := 0; i < 10; i++ {
		wc.callback()
		expect(t, vm.cache.Count(), 1)
		expect(t, vm.scripts.Count(), 1)
		expect(t, vm.version, int32(101+i))
	}
}
func TestVM3(t *testing.T) {
	vm := NewLuaVM(&Binder{}, 1, 1)
	wc := &watcher{callback: vm.Reload}
	vm.watcher = wc
	//检查初始值
	expect(t, vm.version, int32(100))
	expect(t, vm.cache.Count(), 1)
	expect(t, vm.scripts.Count(), 0)

	//检查已加载的脚本数
	err := vm.PreLoad("./test/t1.lua")
	expect(t, err, nil)
	expect(t, vm.scripts.Count(), 1)

	//输入参数错误
	r, m, err := vm.Call("./test/t1.lua", NewContext(""))
	expect(t, err, nil)
	expect(t, len(m), 0)
	expect(t, len(r), 1)
	expect(t, r[0], "hello")

}

func TestVM4(t *testing.T) {
	vm := NewLuaVM(&Binder{}, 1, 1)
	wc := &watcher{callback: vm.Reload}
	vm.watcher = wc
	//检查初始值
	expect(t, vm.version, int32(100))
	expect(t, vm.cache.Count(), 1)
	expect(t, vm.scripts.Count(), 0)

	//检查已加载的脚本数
	_, _, err := vm.Call("./test/t1.lua", NewContext("{}"))
	expect(t, err, nil)
	expect(t, vm.cache.Count(), 1)
	expect(t, vm.scripts.Count(), 0)
}
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
