package lua4go

import (
	"testing"

	"time"

	"github.com/qxnw/lib4go/ut"
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
	vm := NewLuaVM(&Binder{}, 1, 1, time.Hour)
	wc := &watcher{callback: vm.Reload}
	vm.watcher = wc
	//检查初始值
	ut.Expect(t, vm.version, int32(100))
	ut.Expect(t, vm.cache.Count(), 1)
	ut.Expect(t, vm.scripts.Count(), 0)

	//检查已加载的脚本数
	err := vm.PreLoad("./test/t1.lua")
	ut.Expect(t, err, nil)
	ut.Expect(t, vm.scripts.Count(), 1)
	ut.Expect(t, vm.version, int32(100))
	//重复加载脚本
	err = vm.PreLoad("./test/t1.lua")
	ut.Expect(t, err, nil)
	ut.Expect(t, vm.scripts.Count(), 1)
	ut.Expect(t, vm.version, int32(100))

	//加载不存在脚本
	err = vm.PreLoad("./test/t1_not_exist.lua")
	ut.Refute(t, err, nil)
	ut.Expect(t, vm.scripts.Count(), 1)
	ut.Expect(t, vm.version, int32(100))

}

func TestVM2(t *testing.T) {
	vm := NewLuaVM(&Binder{}, 1, 1, time.Hour)
	wc := &watcher{callback: vm.Reload}
	vm.watcher = wc
	//检查初始值
	ut.Expect(t, vm.version, int32(100))
	ut.Expect(t, vm.cache.Count(), 1)
	ut.Expect(t, vm.scripts.Count(), 0)

	//检查已加载的脚本数
	err := vm.PreLoad("./test/t1.lua")
	ut.Expect(t, err, nil)
	ut.Expect(t, vm.scripts.Count(), 1)

	//回调后检查缓存引擎及脚本数
	for i := 0; i < 10; i++ {
		wc.callback()
		ut.Expect(t, vm.cache.Count(), 1)
		ut.Expect(t, vm.scripts.Count(), 1)
		ut.Expect(t, vm.version, int32(101+i))
	}
}
func TestVM3(t *testing.T) {
	vm := NewLuaVM(&Binder{}, 1, 1, time.Hour)
	wc := &watcher{callback: vm.Reload}
	vm.watcher = wc
	//检查初始值
	ut.Expect(t, vm.version, int32(100))
	ut.Expect(t, vm.cache.Count(), 1)
	ut.Expect(t, vm.scripts.Count(), 0)

	//检查已加载的脚本数
	err := vm.PreLoad("./test/t1.lua")
	ut.Expect(t, err, nil)
	ut.Expect(t, vm.scripts.Count(), 1)

	//输入参数错误
	r, m, err := vm.Call("./test/t1.lua", NewContext(""))
	ut.Expect(t, err, nil)
	ut.Expect(t, len(m), 0)
	ut.Expect(t, len(r), 1)
	ut.Expect(t, r[0], "hello")

}

func TestVM4(t *testing.T) {
	vm := NewLuaVM(&Binder{}, 1, 1, time.Hour)
	wc := &watcher{callback: vm.Reload}
	vm.watcher = wc
	//检查初始值
	ut.Expect(t, vm.version, int32(100))
	ut.Expect(t, vm.cache.Count(), 1)
	ut.Expect(t, vm.scripts.Count(), 0)

	//检查已加载的脚本数
	_, _, err := vm.Call("./test/t1.lua", NewContext("{}"))
	ut.Expect(t, err, nil)
	ut.Expect(t, vm.cache.Count(), 1)
	ut.Expect(t, vm.scripts.Count(), 0)
}
