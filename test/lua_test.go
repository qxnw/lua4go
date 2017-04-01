package test

/*
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
func TestTEngineT3(t *testing.T) {
	engine, err := lua4go.NewLuaEngine("./t1.lua", bind.NewDefault())
	ut.Expect(t, err, nil)
	for i := 0; i < 100000; i++ {
		_, _, err := engine.Call(lua4go.NewContext("{}"))
		if ut.ExpectSkip(t, err, nil) {
			return
		}
		if i%100 == 0 {
			t.Logf("%+v", memory.GetInfo().Used)
		}
	}
}
*/
