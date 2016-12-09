package lua4go

//LuaPool lua引擎池
type LuaPool struct {
	binder  *Binder
	minSize int
	maxSize int
}

//NewLuaPool 构建LUA引擎池
func NewLuaPool(binder *Binder, maxSize int, minSize int) *LuaPool {
	pool := &LuaPool{binder: binder, minSize: minSize, maxSize: maxSize}
	return pool
}

//Call 选取最新使用的引擎并根据输入参数执行
func (p *LuaPool) Call(script string, context Context) (result []string, params map[string]string, err error) {

	return

}

//PreLoad 预加载脚本引擎
func (p *LuaPool) PreLoad(script string) (err error) {
	return
}
