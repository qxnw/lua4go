package test

import (
	"testing"

	"fmt"

	"time"

	"github.com/qxnw/lib4go/sysinfo/memory"
	"github.com/qxnw/lib4go/ut"
	"github.com/qxnw/lua4go"
	"github.com/qxnw/lua4go/bind"
)

func BenchmarkEngine1(t *testing.B) {

	t.Logf("%+v", memory.GetInfo().Used)
	for i := 0; i < 1000000; i++ {
		engine, err := lua4go.NewLuaEngine("./t31.lua", bind.NewDefault())
		ut.Expect(t, err, nil)
		//_, _, _ :=
		//m, _, err :=
		engine.Call(lua4go.NewContext("{}"))
		//ut.Expect(t, err, nil)
		//ut.Expect(t, len(m), 1)
		//ut.Expect(t, m[0], "hello")
		//if ut.ExpectSkip(t, err, nil) {
		//return
		//	}
		if i%10000 == 0 {
			fmt.Printf("%+v BenchmarkEngine1:%+v\n", time.Now(), memory.GetInfo().Used)
		}
	}
}
