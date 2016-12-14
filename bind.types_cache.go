package lua4go

import (
	"github.com/qxnw/lib4go/memcache"
	lua "github.com/yuin/gopher-lua"
)

//memcached操作类，用于lua脚本直接调用
//local mem,err=memcached.new("mem")
//if err~=nil then
//	 print(err)
//	 return
//end
//mem:set("key01","value001")   --添加或修改缓存数据，无超时
//mem:set("key02","value002",300) --  添加或修改缓存数据，超时时长为5分钟
//print(mem:get("key01"))  --获取指定key的缓存数据
//mem.del("key01")  ---删除指定key的缓存数据
//mem.delay("key01",300)  ---将key01的超时时长延长为5分钟后

func getMemcachedBinder() *TypeBinder {
	return &TypeBinder{
		Name: "memcached",
		NewFunc: map[string]lua.LGFunction{
			"new": typeNewMemcached,
		},
		Methods: map[string]lua.LGFunction{
			"get":   typeMemcacheGet,
			"add":   typeMemcacheAdd,
			"set":   typeMemcacheSet,
			"delay": typeMemcacheDelay,
			"del":   typeMemcacheDel,
		},
	}
}

// Constructor
func typeNewMemcached(L *lua.LState) int {
	var err error
	ud := L.NewUserData()
	json := L.CheckString(1)
	ud.Value, err = memcache.NewJSON(json)
	if err != nil {
		return pushValues(L, "", err)
	}
	L.SetMetatable(ud, L.GetTypeMetatable("memcached"))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkMemcached(L *lua.LState) *memcache.MemcacheClient {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*memcache.MemcacheClient); ok {
		return v
	}
	L.RaiseError("bad argument  (memcached client expected, got %s)", ud.Type().String())
	return nil
}

func typeMemcacheGet(L *lua.LState) int {
	p := checkMemcached(L)
	key := L.CheckString(2)
	a := p.Get(key)
	return pushValues(L, a)
}
func typeMemcacheAdd(L *lua.LState) int {
	p := checkMemcached(L)
	key := L.CheckString(2)
	value := L.CheckString(3)
	expiresAt := 0
	if L.GetTop() == 4 {
		expiresAt = L.CheckInt(4)
	}
	a := p.Add(key, value, expiresAt)
	return pushValues(L, a)
}
func typeMemcacheSet(L *lua.LState) int {
	p := checkMemcached(L)
	key := L.CheckString(2)
	value := L.CheckString(3)
	expiresAt := 0
	if L.GetTop() == 4 {
		expiresAt = L.CheckInt(4)
	}
	a := p.Set(key, value, expiresAt)
	return pushValues(L, a)
}
func typeMemcacheDel(L *lua.LState) int {
	p := checkMemcached(L)
	key := L.CheckString(2)
	a := p.Delete(key)
	return pushValues(L, a)
}
func typeMemcacheDelay(L *lua.LState) int {
	p := checkMemcached(L)
	key := L.CheckString(2)
	expiresAt := L.CheckInt(3)
	a := p.Delay(key, expiresAt)
	return pushValues(L, a)
}
