package bind

import (
	"github.com/qxnw/lib4go/encoding/base64"
	"github.com/qxnw/lib4go/security/aes"
	"github.com/qxnw/lib4go/security/des"
	"github.com/qxnw/lib4go/security/md5"
	"github.com/qxnw/lib4go/security/rsa"
	"github.com/qxnw/lib4go/security/sha1"
	"github.com/qxnw/lib4go/security/sha256"
	"github.com/yuin/gopher-lua"
)

func moduleMd5Encrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	return pushValues(ls, md5.Encrypt(input))
}
func moduleDESEncrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	key := ls.CheckString(2)
	r, e := des.Encrypt(input, key)
	if e != nil {
		return pushValues(ls, r, e)
	}
	return pushValues(ls, r)
}

func moduleDESDecrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	key := ls.CheckString(2)
	r, e := des.Decrypt(input, key)
	if e != nil {
		return pushValues(ls, r, e)
	}
	return pushValues(ls, r)
}

func moduleQXDESEncrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	key := md5.Encrypt("_QX_ARS_KEY_&" + ls.CheckString(2))
	r, e := des.Encrypt(input, key[0:8])
	if e != nil {
		return pushValues(ls, r, e)
	}
	return pushValues(ls, r)
}
func moduleQXDESDecrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	key := md5.Encrypt("_QX_ARS_KEY_&" + ls.CheckString(2))
	r, e := des.Decrypt(input, key[0:8])
	if e != nil {
		return pushValues(ls, r, e)
	}
	return pushValues(ls, r)
}
func moduleAESEncrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	key := ls.CheckString(2)
	r, e := aes.Encrypt(input, key)
	return pushValues(ls, r, e)
}
func moduleAESDecrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	key := ls.CheckString(2)
	r, e := aes.Decrypt(input, key)
	return pushValues(ls, r, e)
}
func moduleBase64Encode(ls *lua.LState) int {
	input := ls.CheckString(1)
	return pushValues(ls, base64.Encode(input))
}
func moduleBase64EncodeBytes(ls *lua.LState) int {
	input := ls.CheckUserData(1)
	data := input.Value.([]byte)
	return pushValues(ls, base64.EncodeBytes(data))
}

func moduleBase64Decode(ls *lua.LState) int {
	input := ls.CheckString(1)
	r, e := base64.Decode(input)
	return pushValues(ls, r, e)
}
func moduleBase64DecodeBytes(ls *lua.LState) int {
	input := ls.CheckString(1)
	r, e := base64.DecodeBytes(input)
	return pushValues(ls, r, e)
}
func moduleSha256Encrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	return pushValues(ls, sha256.Encrypt(input))
}
func moduleSha1Encrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	return pushValues(ls, sha1.Encrypt(input))
}
func moduleRsaEncrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	publicKey := ls.CheckString(2)
	data, err := rsa.Encrypt(input, publicKey)
	return pushValues(ls, data, err)
}
func moduleRsaDecrypt(ls *lua.LState) int {
	input := ls.CheckString(1)
	privateKey := ls.CheckString(2)
	data, err := rsa.Decrypt(input, privateKey)
	return pushValues(ls, data, err)
}
func moduleRsaMakeSign(ls *lua.LState) int {
	input := ls.CheckString(1)
	privateKey := ls.CheckString(2)
	mode := "sha1"
	if ls.GetTop() > 2 {
		mode = ls.CheckString(3)
	}
	data, err := rsa.Sign(input, privateKey, mode)
	return pushValues(ls, data, err)
}
func moduleRsaVerify(ls *lua.LState) int {
	src := ls.CheckString(1)
	sign := ls.CheckString(2)
	pubkey := ls.CheckString(3)
	mode := "sha1"
	if ls.GetTop() > 3 {
		mode = ls.CheckString(4)
	}
	data, err := rsa.Verify(src, sign, pubkey, mode)
	return pushValues(ls, data, err)
}
