package bind

import (
	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/mq"
	lua "github.com/yuin/gopher-lua"
)

//mq操作类，用于lua脚本直接调用
//local producer=mq.new("mq")  --根据zk配置名称初始化mq
//producer:send(queue,content,timeout) --发送消息

var mqProducerCache = cmap.New()

func getMQTypeBinder() *TypeBinder {
	return &TypeBinder{
		Name: "mq",
		NewFunc: map[string]lua.LGFunction{
			"new": typeMQProducerType,
		},
		Methods: map[string]lua.LGFunction{
			"send": typeMQProducerSend,
		},
	}
}

// Constructor
func typeMQProducerType(L *lua.LState) int {
	var err error
	ud := L.NewUserData()
	config := L.CheckString(1)
	_, producer, err := mqProducerCache.SetIfAbsentCb(config, func(p ...interface{}) (interface{}, error) {
		config := p[0].(string)
		return mq.NewStompProducerJSON(config)
	}, config)
	if err != nil {
		return pushValues(L, nil, err)
	}

	ud.Value = producer
	L.SetMetatable(ud, L.GetTypeMetatable("mqproducer"))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkMQProducerType(L *lua.LState) *mq.StompProducer {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*mq.StompProducer); ok {
		return v
	}
	L.RaiseError("bad argument  (MQProducer expected, got %s)", ud.Type().String())
	return nil
}

func typeMQProducerSend(L *lua.LState) int {
	p := checkMQProducerType(L)
	queue := L.CheckString(2)
	content := L.CheckString(3)
	timeout := L.CheckInt(4)
	a := p.Send(queue, content, timeout)
	return pushValues(L, a)
}
