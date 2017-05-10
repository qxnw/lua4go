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

var (
	ErrVMClose = errors.New("vm is closing")
)

type option struct {
	watchScriptSpan time.Duration
	minSize         int
	maxSize         int
	timeout         time.Duration
}

//Option 配置选项
type Option func(*option)

//WithWatchScript 定时监控脚本
func WithWatchScript(t time.Duration) Option {
	return func(o *option) {
		o.watchScriptSpan = t
	}
}

//WithMinSize 设置最小缓存数
func WithMinSize(t int) Option {
	return func(o *option) {
		o.minSize = t
	}
}

//WithMaxSize 设置最大缓存数
func WithMaxSize(t int) Option {
	return func(o *option) {
		o.maxSize = t
	}
}

//WithTimeout 设置超时时长
func WithTimeout(t time.Duration) Option {
	return func(o *option) {
		o.timeout = t
	}
}

//LuaVM lua虚拟机
type LuaVM struct {
	version int32
	binder  IBinder
	watcher fileWatcher
	cache   cmap.ConcurrentMap
	scripts cmap.ConcurrentMap
	done    bool
	lk      sync.Mutex
	*option
}

//NewLuaVM   构建LUA对象池
func NewLuaVM(binder IBinder, opts ...Option) *LuaVM {
	vm := &LuaVM{binder: binder, version: 99}
	vm.option = &option{minSize: 1, maxSize: 99, timeout: time.Second * 300}
	for _, opt := range opts {
		opt(vm.option)
	}
	if vm.watchScriptSpan > 0 {
		vm.watcher = file.NewDirWatcher(vm.Reload, vm.watchScriptSpan)
	}

	vm.cache = cmap.New()
	vm.scripts = cmap.New()
	vm.cache.SetIfAbsentCb(string(vm.version+1), vm.createNewPool)
	atomic.AddInt32(&vm.version, 1)
	return vm
}

//Call 选取最新的脚本引擎执行当前脚本
func (vm *LuaVM) Call(script string, input *Context) (result []string, params map[string]string, err error) {
	if vm.done {
		err = ErrVMClose
		return
	}
	pl, b := vm.cache.Get(string(vm.version))
	if !b {
		err = ErrVMClose
		return
	}
	defer vm.addWatch(script)
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
		err = ErrVMClose
		return
	}
	pl, _ := vm.cache.Get(string(vm.version))
	_, err = pl.(*LuaPool).PreLoad(script)
	if err != nil {
		return
	}
	defer vm.addWatch(script)
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
	pl := NewLuaPool(vm.binder, vm.minSize, vm.maxSize, vm.timeout)
	vm.scripts.IterCb(func(k string, v interface{}) bool {
		pl.PreLoad(k)
		return true
	})
	return pl, nil
}
func (vm *LuaVM) addWatch(path string) {
	if vm.watcher == nil {
		return
	}
	vm.watcher.Append(path)
}
