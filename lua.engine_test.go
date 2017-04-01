package lua4go

import (
	"fmt"
	"testing"

	"github.com/qxnw/lib4go/ut"
)

//测试引擎创建，基本的脚本调用，返回值转换，输入参数转换
func TestEngine1(t *testing.T) {
	engine, err := NewLuaEngine("./test/t1.lua", &Binder{})
	expect(t, err, nil)
	expect(t, engine.script, "./test/t1.lua")
	refute(t, engine.state, nil)
}
func TestEngine2(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t1.lua", &Binder{})
	r, m, err := engine.Call(NewContext("{}"))
	expect(t, err, nil)
	expect(t, len(m), 0)
	expect(t, len(r), 1)
	expect(t, r[0], "hello")

}
func TestEngine20(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t1.lua", &Binder{})
	r, m, err := engine.Call(NewContext("{}"))
	r, m, err = engine.Call(NewContext(""))
	expect(t, err, nil)
	expect(t, len(m), 0)
	expect(t, len(r), 1)
	expect(t, r[0], "hello")
}
func TestEngine3(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t2.lua", &Binder{})
	//输入正确参数，返回两个值
	r, m, err := engine.Call(NewContext("{}"))
	expect(t, err, nil)
	expect(t, len(m), 0)
	expect(t, len(r), 1)
	expect(t, r[0], "hello")
}
func TestEngine4(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t3.lua", &Binder{})
	//输入正确参数，返回3个值，但只能处理2个值
	r, m, err := engine.Call(NewContext("{}"))
	expect(t, err, nil)
	expect(t, len(m), 0)
	expect(t, len(r), 1)
	expect(t, r[0], "")
}
func TestEngine5(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t4.lua", &Binder{})
	//输入正确参数，返回3个值，但只能处理2个值
	r, m, err := engine.Call(NewContext("{}"))
	expect(t, err, nil)
	expect(t, len(m), 0)
	expect(t, len(r), 1)
	expect(t, r[0], "{}")
}
func TestEngine6(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t5.lua", &Binder{})
	//输入正确参数，返回3个值，但只能处理2个值
	r, m, err := engine.Call(NewContext("{}"))
	expect(t, err, nil)
	expect(t, len(m), 0)
	expect(t, len(r), 1)
	expect(t, r[0], `{"id":"1"}`)
}
func TestEngine7(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t6.lua", &Binder{})
	//输入正确参数，返回3个值，但只能处理2个值
	r, m, err := engine.Call(NewContext("{}"))
	expect(t, err, nil)
	expect(t, len(m), 0)
	expect(t, len(r), 1)
	expect(t, r[0], `{"id":"1","value":{"id":"2"}}`)
}
func TestEngine8(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t7.lua", &Binder{})
	//输入正确参数，返回3个值，但只能处理2个值
	_, _, err := engine.Call(NewContext("{}"))
	refute(t, err, nil)
}
func TestEngine81(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t8.lua", &Binder{})
	//输入正确参数，返回3个值，但只能处理2个值
	r, _, err := engine.Call(NewContext(`{"x":100,"y":0}`))
	fmt.Println(err)
	ut.Refute(t, err, nil)
	ut.Expect(t, len(r), 0)

	r, _, err = engine.Call(NewContext(`{"x":100,"y":1}`))
	ut.Expect(t, err, nil)
	ut.Expect(t, len(r), 1)
	ut.Expect(t, r[0], "100")
}

func TestEngine9(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t1.lua", &Binder{})
	//输入正确参数，返回3个值，但只能处理2个值
	_, _, err := engine.Call(NewContext("{id=2}"))
	refute(t, err, nil)
}
func TestEngine10(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t10.lua", &Binder{})
	//输入正确参数，返回3个值，但只能处理2个值
	r, m, err := engine.Call(NewContext(`{"id":2}`))
	expect(t, err, nil)
	expect(t, len(m), 0)
	expect(t, len(r), 1)
	expect(t, r[0], `2`)
}
func TestEngine11(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t10.lua", &Binder{})
	//输入正确参数，返回3个值，但只能处理2个值
	r, m, err := engine.Call(NewContext(`{"id":"1","value":{"id":"2"}}`))
	expect(t, err, nil)
	expect(t, len(m), 0)
	expect(t, len(r), 1)
	expect(t, r[0], `1`)
}
func TestEngine12(t *testing.T) {
	engine, _ := NewLuaEngine("./test/t11.lua", &Binder{})
	//输入正确参数，返回3个值，但只能处理2个值
	r, m, err := engine.Call(NewContext(`{"id":"1","value":{"id":"2"}}`))
	expect(t, err, nil)
	expect(t, len(m), 0)
	expect(t, len(r), 1)
	expect(t, r[0], `2`)
}
