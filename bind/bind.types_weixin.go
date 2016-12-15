package bind

import (
	"github.com/qxnw/lib4go/weixin"
	lua "github.com/yuin/gopher-lua"
)

//weixin操作类，用于lua脚本直接调用
//local wx=weixin.new(appid, token, encodingAESKey)  --根据appid, token, encodingAESKey构建weixin操作类
//local content,err=wx.decrypt(encryptedContent) --解密
//local content,err=wx.encrypt(formUserName, toUsername, content, monce, timestamp) --加密
//local sign=wx.makeSign(timestamp, monce, msg) //构建签名串

func getWeixinTypeBinder() *TypeBinder {
	return &TypeBinder{
		Name: "weixin",
		NewFunc: map[string]lua.LGFunction{
			"new": typeWeixinType,
		},
		Methods: map[string]lua.LGFunction{
			"decrypt":   typeWeixinDecrypt,
			"encrypt":   typeWeixinEncrypt,
			"make_sign": typeWeixinMakeSign,
		},
	}
}

// Constructor
func typeWeixinType(L *lua.LState) int {
	var err error
	ud := L.NewUserData()
	appid := L.CheckString(1)
	token := L.CheckString(2)
	encodingAESKey := L.CheckString(3)
	ud.Value, err = weixin.NewWechat(appid, token, encodingAESKey)
	if err != nil {
		return pushValues(L, "", err)
	}
	L.SetMetatable(ud, L.GetTypeMetatable("weixin"))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkWeixinType(L *lua.LState) weixin.Wechat {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(weixin.Wechat); ok {
		return v
	}
	L.RaiseError("bad argument  (http client expected, got %s)", ud.Type().String())
	return weixin.Wechat{}
}

func typeWeixinDecrypt(L *lua.LState) int {
	p := checkWeixinType(L)
	content := L.CheckString(2)
	a, b := p.Decrypt(content)
	return pushValues(L, a, b)
}
func typeWeixinEncrypt(L *lua.LState) int {
	p := checkWeixinType(L)
	formUserName := L.CheckString(2)
	toUsername := L.CheckString(3)
	content := L.CheckString(4)
	monce := L.CheckString(5)
	timestamp := L.CheckString(6)
	a, b := p.Encrypt(formUserName, toUsername, content, monce, timestamp)
	return pushValues(L, a, b)
}
func typeWeixinMakeSign(L *lua.LState) int {
	p := checkWeixinType(L)
	timestamp := L.CheckString(2)
	monce := L.CheckString(3)
	msg := L.CheckString(4)
	a := p.MakeSign(timestamp, monce, msg)
	return pushValues(L, a)
}
