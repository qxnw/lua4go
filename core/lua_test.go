package core

import (
	"strings"
	"testing"

	"github.com/arsgo/lib4go/security/md5"
)

var luaPool *LuaPool
var min int
var max int

func init() {
	min = 100
	max = 1000
	luaPool = NewLuaPool()
}
func TestInit(t *testing.T) {
	if luaPool.PreLoad("./t1.lua", min, max) != nil {
		t.Error("luapool init error")
	}
	if luaPool.PreLoad("./t2.lua", min, max) != nil {
		t.Error("luapool init error")
	}
}

/*
func TestBenchCall(t *testing.T) {
	time.Sleep(time.Second * 2)
	ch := make(chan int, max)
	close := make(chan int, 1)
	var index int32
	var concurrent int32
	concurrent = 100000
	groupName := "./t2.lua"

	for i := 0; i < min; i++ {
		ch <- i
		go func() {
			for {
				if atomic.LoadInt32(&index) >= concurrent {
					close <- 1
					break
				}
				<-ch
				values, _, err := luaPool.Call(groupName, "", "{}", "123456")
				if err != nil {
					t.Error(err.Error())
				} else {
					if len(values) != 1 {
						t.Error("return values len error")
					}
				}
				atomic.AddInt32(&index, 1)
				ch <- 1
			}

		}()
	}
	<-close

}
*/
func TestLua(t *testing.T) {

	values, _, err := luaPool.Call("./t2.lua", "", "{}", "123456")
	if err != nil {
		t.Error(err)
	}
	if len(values) != 2 {
		t.Error("return values len error", len(values))
		t.SkipNow()
	}

	if !strings.EqualFold(md5.Encrypt("123456"), values[0]) {
		t.Errorf("return values is error %s,expect:%s", values[0], md5.Encrypt("123456"))
	}
	if !strings.EqualFold("123456", values[1]) {
		t.Errorf("return values is error %s,expect:123456", values[1])
	}

}
