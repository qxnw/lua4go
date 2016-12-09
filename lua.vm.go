package lua4go

import (
	"errors"
	"sync"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/file"
)

//LuaVM lua虚拟机
type LuaVM struct {
	version int32
	binder  *Binder
	watcher *file.DirWatcher
	cache   cmap.ConcurrentMap
	minSize int
	maxSize int
	lk      sync.Mutex
}

//NewLuaVM   构建LUA对象池
func NewLuaVM(binder *Binder) *LuaVM {
	vm := &LuaVM{binder: binder, version: 0}
	vm.watcher = file.NewDirWatcher(vm.Reload)
	vm.cache = cmap.New()
	vm.cache.SetIfAbsentCb(string(vm.version+1), vm.createVM, vm.version, vm.version+1)
	return vm
}

//SetPoolSize 设置连接池大小
func (p *LuaVM) SetPoolSize(minSize int, maxSize int) {
	p.minSize = minSize
	p.maxSize = maxSize
}

//Call 选取最新的脚本引擎执行当前脚本
func (p *LuaVM) Call(script string, input Context) (result []string, outparams map[string]string, err error) {
	return
}

//Reload 重新加载所有引擎
func (p *LuaVM) Reload() {
}

//PreLoad 预加载脚本
func (p *LuaVM) PreLoad(script string, minSize int, maxSize int) (err error) {
	p.watcher.Append(script)
	return
}

//Close 关闭引擎
func (p *LuaVM) Close() {
}

func (p *LuaVM) createVM(args ...interface{}) (interface{}, error) {
	p.lk.Lock()
	defer p.lk.Unlock()

	return nil, errors.New("创建失败，版本错误")
}
