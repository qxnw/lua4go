package lua4go

import (
	"strings"

	"github.com/qxnw/lib4go/net/http"
	lua "github.com/yuin/gopher-lua"
)

//http操作类，用于lua脚本直接调用
//local client=http.new() ---构建普通的http client
//local client=http.new("http://192.168.0.1:8080") ---使用http代理服务器构建http client
//local client=http.new("/f1.cert","/f2.key","/f3.ca") ---使用加密证书构建http client
//client:get("http://localhost:1016/system/md5") ---发送http.get请求，默认编码为utf-8
//client:get("http://localhost:1016/system/md5","gbk") ---发送http.get请求，编码为gbk
//client:post("http://localhost:1016/system/md5","id=1") ---发送http.post请求，默认编码为utf-8
//client:post("http://localhost:1016/system/md5","id=1","gbk") ---发送http.post请求，编码为gbk
//client:request("put","http://localhost:1016/system/md5","id=1","gbk",{charset="gbk"}) ---发送http.put请求，编码为gbk
//client:request("delete","http://localhost:1016/system/md5","id=1","gbk",{charset="gbk"}) ---发送http.delete请求，编码为gbk

func getHTTPClientTypeBinder() *TypeBinder {
	return &TypeBinder{
		Name: "http",
		NewFunc: map[string]lua.LGFunction{
			"new": typeNewHTTP,
		},
		Methods: map[string]lua.LGFunction{
			"get":      typeDoHttpGet,
			"post":     typeDoHttpPost,
			"request":  typeDoHttpRequest,
			"download": typeDoHttpDownload,
			"save":     typeDoHttpSave,
		},
	}
}

// Constructor
func typeNewHTTP(L *lua.LState) int {
	var client *http.HTTPClient
	var err error
	count := L.GetTop()
	switch count {
	case 0:
		client = http.NewHTTPClient()
	case 1:
		proxy := L.CheckString(1)
		if strings.HasPrefix(strings.ToLower(proxy), "http:") ||
			strings.HasPrefix(strings.ToLower(proxy), "tcp:") {
			client = http.NewHTTPClientProxy(proxy)
		} else {
			client, err = http.NewHTTPClientCert1(proxy)
			if err != nil {
				return pushValues(L, "", err)
			}
		}
	case 3:
		certFile := L.CheckString(1)
		keyFile := L.CheckString(2)
		caFile := L.CheckString(3)
		client, err = http.NewHTTPClientCert(certFile, keyFile, caFile)
		if err != nil {
			return pushValues(L, "", err)
		}
	default:
		L.RaiseError("bad argument (len must less than 4)")
	}
	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable("HTTPClient"))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkHTTP(L *lua.LState) *http.HTTPClient {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*http.HTTPClient); ok {
		return v
	}
	L.RaiseError("bad argument  (http client expected, got %s)", ud.Type().String())
	return nil
}

func typeDoHttpRequest(L *lua.LState) int {
	p := checkHTTP(L)
	method := L.CheckString(2)
	url := L.CheckString(3)
	params := L.CheckString(4)
	encoding := L.CheckString(5)
	header := L.CheckTable(6)
	a, b, c := p.Request(method, url, params, encoding, getMapParams(header))
	return pushValues(L, a, b, c)
}

func typeDoHttpDownload(L *lua.LState) int {
	p := checkHTTP(L)
	method := L.CheckString(2)
	url := L.CheckString(3)
	params := L.CheckString(4)
	header := L.CheckTable(5)
	a, b, c := p.Download(method, url, params, getMapParams(header))
	return pushValues(L, a, b, c)
}

func typeDoHttpSave(L *lua.LState) int {
	p := checkHTTP(L)
	method := L.CheckString(2)
	url := L.CheckString(3)
	params := L.CheckString(4)
	header := L.CheckTable(5)
	path := L.CheckString(6)
	a, b := p.Save(method, url, params, getMapParams(header), path)
	return pushValues(L, a, b)
}

func typeDoHttpGet(L *lua.LState) int {
	p := checkHTTP(L)
	url := L.CheckString(2)
	parms := getStringParams(L, 3)
	a, b, c := p.Get(url, parms...)
	return pushValues(L, a, b, c)
}

func typeDoHttpPost(L *lua.LState) int {
	p := checkHTTP(L)
	url := L.CheckString(2)
	data := L.CheckString(3)
	parms := getStringParams(L, 4)
	a, b, c := p.Post(url, data, parms...)
	return pushValues(L, a, b, c)
}
