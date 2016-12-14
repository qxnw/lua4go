package lua4go

import (
	"github.com/arsgo/lib4go/encoding"
	"github.com/qxnw/lib4go/encoding/html"
	"github.com/qxnw/lib4go/encoding/url"
	"github.com/yuin/gopher-lua"
)

func moduleEncodingConvert(ls *lua.LState) int {
	input := ls.CheckString(1)
	chaset := ls.CheckString(2)
	result := encoding.Convert([]byte(input), chaset)
	return pushValues(ls, result)
}
func moduleUnicodeEncode(ls *lua.LState) int {
	input := ls.CheckString(1)
	result := encoding.UnicodeEncode(input)
	return pushValues(ls, result)
}
func moduleUnicodeDecode(ls *lua.LState) int {
	input := ls.CheckString(1)
	result := encoding.UnicodeDecode(input)
	return pushValues(ls, result)
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
