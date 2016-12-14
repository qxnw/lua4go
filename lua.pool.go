package lua4go

import (
	"sync/atomic"
	"time"

	"errors"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/pool"
)

//IdleTimeout 最长空闲时间
var IdleTimeout = time.Second * 30

//LuaPool lua引擎池
type LuaPool struct {
	binder  *Binder
	minSize int
	maxSize int
	using   int32
	isClose bool
	cache   cmap.ConcurrentMap
}

//NewLuaPool 构建LUA引擎池
func NewLuaPool(binder *Binder, minSize int, maxSize int) *LuaPool {
	pool := &LuaPool{binder: binder, minSize: minSize, maxSize: maxSize, isClose: false}
	pool.cache = cmap.New()
	return pool
}

//Call 选取最新使用的引擎并根据输入参数执行
func (p *LuaPool) Call(script string, context *Context) (result []string, params map[string]string, err error) {
	if p.isClose {
		err = errors.New("脚本引擎已关闭")
		return
	}
	atomic.AddInt32(&p.using, 1)
	defer atomic.AddInt32(&p.using, -1)
	pl, err := p.PreLoad(script)
	if err != nil {
		return
	}
	engine, err := pl.Get()
	if err != nil {
		return
	}
	defer pl.Put(engine)
	return engine.(*LuaEngine).Call(context)
}

//PreLoad 预加载脚本引擎
func (p *LuaPool) PreLoad(script string) (pl pool.IPool, err error) {
	if p.isClose {
		err = errors.New("脚本引擎已经关闭")
		return
	}
	if _, obj, err := p.cache.SetIfAbsentCb(script, func(input ...interface{}) (interface{}, error) {
		script := input[0].(string)
		return pool.New(&pool.PoolConfigOptions{
			InitialCap:  p.minSize,
			MaxCap:      p.maxSize,
			IdleTimeout: IdleTimeout,
			Factory: func() (interface{}, error) {
				return NewLuaEngine(script, p.binder)
			},
			Close: func(v interface{}) error {
				engine := v.(*LuaEngine)
				engine.Close()
				return nil
			},
		})

	}, script); err == nil {
		pl = obj.(pool.IPool)
	}
	return
}

//Close 关闭所有连接池
func (p *LuaPool) Close() {
	p.isClose = true
	go func() {
		tk := time.NewTicker(time.Second * 5)
	START:
		for {
			select {
			case <-tk.C:
				if p.using == 0 {
					p.cache.RemoveIterCb(func(key string, pl interface{}) bool {
						pl.(pool.IPool).Release()
						return true
					})
					tk.Stop()
					break START
				}
			}
		}
	}()
}
