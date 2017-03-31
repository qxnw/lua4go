package bind

import (
	"github.com/qxnw/lib4go/influxdb"
	lua "github.com/yuin/gopher-lua"
)

//influxdb操作类，用于lua脚本直接调用
//local influx,err=influxdb.new("influx")
//if err~=nil then
//	 print(err)
//end
//influx:save("{"id":1,"name":"colin"}")
func getinfluxTypeBinder() *TypeBinder {
	return &TypeBinder{
		Name: "influxdb",
		NewFunc: map[string]lua.LGFunction{
			"new": typeInfluxType,
		},
		Methods: map[string]lua.LGFunction{
			"save": typeInfluxDBSave,
		},
	}
}

// Constructor
func typeInfluxType(L *lua.LState) int {
	var err error
	ud := L.NewUserData()
	name := L.CheckString(1)
	ud.Value, err = influxdb.NewJSON(name)
	if err != nil {
		return pushValues(L, "", err)
	}
	L.SetMetatable(ud, L.GetTypeMetatable("influxdb"))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkInfluxDBType(L *lua.LState) *influxdb.Client {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*influxdb.Client); ok {
		return v
	}
	L.RaiseError("bad argument  (influxdb.InfluxDB expected, got %s)", ud.Type().String())
	return nil
}
func typeInfluxDBSave(L *lua.LState) int {
	p := checkInfluxDBType(L)
	data := L.CheckString(2)
	database := L.CheckString(3)

	retentionPolicy := ""
	precision := ""
	writeConsistency := ""
	if L.GetTop() > 3 {
		retentionPolicy = L.CheckString(4)
	}
	if L.GetTop() > 4 {
		precision = L.CheckString(5)
	}
	if L.GetTop() > 5 {
		writeConsistency = L.CheckString(6)
	}
	_, err := p.WriteLineProtocol(data, database, retentionPolicy, precision, writeConsistency)
	if err != nil {
		return pushValues(L, "", err)
	}
	return pushValues(L)
}
