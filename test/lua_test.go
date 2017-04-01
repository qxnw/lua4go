package test

import (
	"testing"

	"github.com/qxnw/lib4go/ut"
	"github.com/qxnw/lua4go"
	"github.com/qxnw/lua4go/bind"
)

func TestTEngineT1(t *testing.T) {
	engine, err := lua4go.NewLuaEngine("./t30.lua", bind.NewDefault())
	ut.Expect(t, err, nil)
	r, m, err := engine.Call(lua4go.NewContext("{}"))
	ut.Expect(t, err, nil)
	ut.Expect(t, len(r), 1)
	ut.Expect(t, len(m), 0)
	ut.Expect(t, r[0], "m")
}
func TestTEngineT2(t *testing.T) {
	engine, err := lua4go.NewLuaEngine("./t31.lua", bind.NewDefault())
	ut.Expect(t, err, nil)
	_, _, err = engine.Call(lua4go.NewContext("{}"))
	ut.Refute(t, err, nil)
}
