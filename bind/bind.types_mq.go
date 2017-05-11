package bind

import (
	"fmt"

	"time"

	"github.com/qxnw/lib4go/concurrent/cmap"
	"github.com/qxnw/lib4go/jsons"
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
func typeMQProducerType(ls *lua.LState) int {
	var err error

	name := ls.CheckString(1)
	value, err := getFuncVarGet(ls, "mq", name)
	if err != nil {
		return pushValues(ls, nil, err)
	}
	producer, err := getMQFromCache(ls, value)
	if err != nil {
		return pushValues(ls, nil, err)
	}
	ud := ls.NewUserData()
	ud.Value = producer
	ls.SetMetatable(ud, ls.GetTypeMetatable("mq"))
	ls.Push(ud)
	return 1
}

var mqCache cmap.ConcurrentMap

func init() {
	mqCache = cmap.New()
}

func getMQFromCache(ls *lua.LState, conf string) (mq.MQProducer, error) {
	_, v, err := mqCache.SetIfAbsentCb(conf, func(input ...interface{}) (interface{}, error) {
		config := input[0].(string)
		configMap, err := jsons.Unmarshal([]byte(conf))
		if err != nil {
			return nil, err
		}
		address, ok := configMap["address"]
		if !ok {
			return nil, fmt.Errorf("db配置文件错误，未包含address节点:%s", config)
		}
		opts := make([]mq.Option, 0, 1)
		logger, err := globalGetLogger(ls)
		if err == nil {
			opts = append(opts, mq.WithLogger(logger))
		}
		p, err := mq.NewMQProducer(address.(string), opts...)
		if err != nil {
			return p, err
		}
		err = p.Connect()
		return p, err
	}, conf)
	if err != nil {
		return nil, err
	}
	return v.(mq.MQProducer), nil
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkMQProducerType(L *lua.LState) mq.MQProducer {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(mq.MQProducer); ok {
		return v
	}
	L.RaiseError("bad argument  (MQProducer expected, got %s)", ud.Type().String())
	return nil
}

func typeMQProducerSend(L *lua.LState) int {
	p := checkMQProducerType(L)
	queue := L.CheckString(2)
	input := L.CheckTable(3)
	timeout := 300
	if L.GetTop() > 3 {
		timeout = L.CheckInt(4)
	}
	data := getIMapParams(input)
	buf, err := jsons.Marshal(data)
	if err != nil {
		return pushValues(L, err)
	}
	err = p.Send(queue, string(buf), time.Duration(timeout)*time.Second)
	return pushValues(L, err)
}
