package lua4go

import (
	"sync"
	"sync/atomic"

	"errors"

	"time"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/file"
)

type fileWatcher interface {
	Append(string) error
}

//LuaVM lua虚拟机
type LuaVM struct {
	version int32
	binder  IBinder
	watcher fileWatcher
	cache   cmap.ConcurrentMap
	scripts cmap.ConcurrentMap
	done    bool
	minSize int
	maxSize int
	lk      sync.Mutex
}

//NewLuaVM   构建LUA对象池
func NewLuaVM(binder IBinder, minSize int, maxSize int) *LuaVM {
	vm := &LuaVM{binder: binder, version: 99, minSize: minSize, maxSize: maxSize}
	vm.watcher = file.NewDirWatcher(vm.Reload, time.Second*5)
	vm.cache = cmap.New()
	vm.scripts = cmap.New()
	vm.cache.SetIfAbsentCb(string(vm.version+1), vm.createNewPool)
	atomic.AddInt32(&vm.version, 1)
	return vm
}

//Call 选取最新的脚本引擎执行当前脚本
func (vm *LuaVM) Call(script string, input *Context) (result []string, params map[string]string, err error) {
	if vm.done {
		err = errors.New("虚拟机已关闭")
		return
	}
	pl, b := vm.cache.Get(string(vm.version))
	if !b {
		err = errors.New("内部错误未找到引擎")
		return
	}
	defer vm.watcher.Append(script)
	return pl.(*LuaPool).Call(script, input)
}

//Reload 重新加载所有引擎
func (vm *LuaVM) Reload() {
	if vm.done {
		return
	}
	vm.lk.Lock()
	defer vm.lk.Unlock()
	oldVersion := string(vm.version)
	oldPool, _ := vm.cache.Get(oldVersion)
	ok, _, _ := vm.cache.SetIfAbsentCb(string(vm.version+1), vm.createNewPool)
	if ok {
		atomic.AddInt32(&vm.version, 1)
		oldPool.(*LuaPool).Close()
		vm.cache.Remove(oldVersion)
	}
}

//PreLoad 预加载脚本
func (vm *LuaVM) PreLoad(script string) (err error) {
	if vm.done {
		err = errors.New("虚拟机已关闭")
		return
	}
	pl, _ := vm.cache.Get(string(vm.version))
	_, err = pl.(*LuaPool).PreLoad(script)
	if err != nil {
		return
	}
	vm.watcher.Append(script)
	vm.scripts.SetIfAbsent(script, script)
	return nil
}

//Close 关闭引擎
func (vm *LuaVM) Close() {
	vm.done = true
	vm.cache.RemoveIterCb(func(key string, p interface{}) bool {
		p.(*LuaPool).Close()
		return true
	})
}

func (vm *LuaVM) createNewPool(args ...interface{}) (p interface{}, er error) {
	pl := NewLuaPool(vm.binder, vm.minSize, vm.maxSize)
	vm.scripts.IterCb(func(k string, v interface{}) bool {
		pl.PreLoad(k)
		return true
	})
	return pl, nil
}
