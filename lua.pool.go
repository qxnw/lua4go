package lua4go

import (
	"sync/atomic"
	"time"

	"errors"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/pool"
)

//IdleTimeout 最长空闲时间

//LuaPool lua引擎池
type LuaPool struct {
	binder  IBinder
	minSize int
	maxSize int
	using   int32
	done    bool
	cache   cmap.ConcurrentMap
	Timeout time.Duration
}

//NewLuaPool 构建LUA引擎池
func NewLuaPool(binder IBinder, minSize int, maxSize int, timeout time.Duration) *LuaPool {
	pool := &LuaPool{binder: binder, minSize: minSize, maxSize: maxSize, Timeout: timeout}
	pool.cache = cmap.New(32)
	return pool
}

//Call 选取最新使用的引擎并根据输入参数执行
func (p *LuaPool) Call(script string, context *Context) (result []string, params map[string]string, err error) {
	if p.done {
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
	eg := engine.(*LuaEngine)
	result, params, err = eg.Call(context)
	pl.Put(engine)
	//eg.Close()
	return

}

//PreLoad 预加载脚本引擎
func (p *LuaPool) PreLoad(script string) (pl pool.IPool, err error) {
	if p.done {
		err = errors.New("脚本引擎已经关闭")
		return
	}
	_, obj, err := p.cache.SetIfAbsentCb(script, func(input ...interface{}) (interface{}, error) {
		script := input[0].(string)
		p, err := pool.New(&pool.PoolConfigOptions{
			InitialCap:  p.minSize,
			MaxCap:      p.maxSize,
			IdleTimeout: p.Timeout,
			Factory: func() (interface{}, error) {
				return NewLuaEngine(script, p.binder)
			},
			Close: func(v interface{}) error {
				engine := v.(*LuaEngine)
				engine.Close()
				return nil
			},
		})
		return p, err

	}, script)

	if err == nil {
		pl = obj.(pool.IPool)
	}
	return pl, err
}

//Close 关闭所有连接池
func (p *LuaPool) Close() {
	p.done = true
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				if p.using == 0 {
					p.cache.RemoveIterCb(func(key string, pl interface{}) bool {
						pl.(pool.IPool).Release()
						return true
					})
					return
				}
			}
		}
	}()
}
