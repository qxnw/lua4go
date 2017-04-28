package bind

import (
	"fmt"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/db"
	"github.com/qxnw/lib4go/jsons"
	lua "github.com/yuin/gopher-lua"
)

//db操作类，用于lua脚本直接调用
//local cdb=db.new("agt_comm") ---根据zk中配置的数据连接名称创建DB实例
//local result,err=cdb:query("select 1 from dual where id=@id",{id=1})
//local result,err=cdb:execute("update tb1 set time=sysdate where id=@id",{id=1})
//local trans=cdb:begin()
//local result,err=trans:cdb:execute("update tb1 set time=sysdate where id=@id",{id=1})
//trans:rollback()
//trans.commit()

func getdbTransTypeBinder() *TypeBinder {
	return &TypeBinder{
		Name:    "dbtrans",
		NewFunc: map[string]lua.LGFunction{},
		Methods: map[string]lua.LGFunction{
			"execute":  typeDBTransExecute,
			"query":    typeDBTransQuery,
			"scalar":   typeDBTransScalar,
			"commit":   typeDBTransCommit,
			"rollback": typeDBTransRollback,
		},
	}
}
func getdbTypeBinder() *TypeBinder {
	return &TypeBinder{
		Name: "db",
		NewFunc: map[string]lua.LGFunction{
			"new": typeDBType,
		},
		Methods: map[string]lua.LGFunction{
			"execute": typeDBExecute,
			"begin":   typeDBBegin,
			"query":   typeDBQuery,
			"scalar":  typeDBScalar,
		},
	}
}

// Constructor
func typeDBType(ls *lua.LState) int {
	name := ls.CheckString(1)
	value, err := getFuncVarGet(ls, "db", name)
	if err != nil {
		return pushValues(ls, "", err)
	}
	ud := ls.NewUserData()
	ud.Value, err = getDBFromCache(value)
	if err != nil {
		return pushValues(ls, "", err)
	}
	ls.SetMetatable(ud, ls.GetTypeMetatable("db"))
	ls.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkDBType(L *lua.LState) *db.DB {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*db.DB); ok {
		return v
	}
	L.RaiseError("bad argument  (db expected, got %s)", ud.Type().String())
	return nil
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkDBTransType(L *lua.LState) *db.DBTrans {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*db.DBTrans); ok {
		return v
	}
	L.RaiseError("bad argument  (db trans expected, got %s)", ud.Type().String())
	return nil
}

func typeDBExecute(L *lua.LState) int {
	p := checkDBType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	row, query, args, err := p.Execute(query, getIMapParams(input))
	return pushValues(L, row, err, query, toLuaTable(L, args))
}

func typeDBQuery(L *lua.LState) int {
	p := checkDBType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	data, query, args, err := p.Query(query, getIMapParams(input))
	if err != nil {
		return pushValues(L, "", err)
	}
	return pushValues(L, data, err, query, toLuaTable(L, args))
}
func typeDBScalar(L *lua.LState) int {

	p := checkDBType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	data, query, args, err := p.Scalar(query, getIMapParams(input))
	return pushValues(L, data, err, query, toLuaTable(L, args))
}
func typeDBBegin(L *lua.LState) int {
	p := checkDBType(L)
	ts, err := p.Begin()
	if err != nil {
		return pushValues(L, "", err)
	}
	ud := L.NewUserData()
	ud.Value = ts
	L.SetMetatable(ud, L.GetTypeMetatable("dbtrans"))
	L.Push(ud)
	return 1
}

func typeDBTransExecute(L *lua.LState) int {
	p := checkDBTransType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	row, query, args, err := p.Execute(query, getIMapParams(input))
	return pushValues(L, row, err, query, toLuaTable(L, args))
}

func typeDBTransQuery(L *lua.LState) int {
	p := checkDBTransType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	data, query, args, err := p.Query(query, getIMapParams(input))
	return pushValues(L, data, err, query, toLuaTable(L, args))
}
func typeDBTransScalar(L *lua.LState) int {
	p := checkDBTransType(L)
	query := L.CheckString(2)
	input := L.CheckTable(3)
	data, query, args, err := p.Scalar(query, getIMapParams(input))
	return pushValues(L, data, err, query, toLuaTable(L, args))
}
func typeDBTransCommit(L *lua.LState) int {
	p := checkDBTransType(L)
	a := p.Commit()
	return pushValues(L, a)
}
func typeDBTransRollback(L *lua.LState) int {
	p := checkDBTransType(L)
	a := p.Rollback()
	return pushValues(L, a)
}

var dbCache cmap.ConcurrentMap

func init() {
	dbCache = cmap.New()
}

func getDBFromCache(conf string) (*db.DB, error) {
	_, v, err := dbCache.SetIfAbsentCb(conf, func(input ...interface{}) (interface{}, error) {
		config := input[0].(string)
		configMap, err := jsons.Unmarshal([]byte(conf))
		if err != nil {
			return nil, err
		}
		provider, ok := configMap["provider"]
		if !ok {
			return nil, fmt.Errorf("db配置文件错误，未包含provider节点:%s", config)
		}
		connString, ok := configMap["connString"]
		if !ok {
			return nil, fmt.Errorf("db配置文件错误，未包含connString节点:%s", config)
		}
		return db.NewDB(provider.(string), connString.(string), 1, 5)

	}, conf)
	if err != nil {
		return nil, err
	}
	return v.(*db.DB), nil
}
