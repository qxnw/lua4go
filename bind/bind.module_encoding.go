package bind

import (
	"github.com/qxnw/lib4go/encoding"
	"github.com/qxnw/lib4go/encoding/html"
	"github.com/qxnw/lib4go/encoding/url"
	"github.com/yuin/gopher-lua"
)

func moduleEncodingConvert(ls *lua.LState) int {
	input := ls.CheckString(1)
	chaset := ls.CheckString(2)
	result, err := encoding.Convert([]byte(input), chaset)
	return pushValues(ls, result, err)
}

func moduleURLEncode(ls *lua.LState) int {
	input := ls.CheckString(1)
	result := url.Encode(input)
	return pushValues(ls, result)
}
func moduleURLDecode(ls *lua.LState) int {
	input := ls.CheckString(1)
	result, err := url.Decode(input)
	return pushValues(ls, result, err)
}
func moduleHTMLEncode(ls *lua.LState) int {
	input := ls.CheckString(1)
	result := html.Encode(input)
	return pushValues(ls, result)
}
func moduleHTMLDecode(ls *lua.LState) int {
	input := ls.CheckString(1)
	result := html.Decode(input)
	return pushValues(ls, result)
}
