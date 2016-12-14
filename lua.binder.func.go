package lua4go

import "github.com/yuin/gopher-lua"

func getGlobal() (r map[string]lua.LGFunction) {
	r = map[string]lua.LGFunction{
		"print":  globalInfo,
		"printf": globalInfof,
		"info":   globalInfo,
		"infof":  globalInfof,
		"error":  globalError,
		"errorf": globalErrorf,
		"sleep":  globalSleep,
		"guid":   globalGUID,
	}
	return
}

func getModules() (r map[string]map[string]lua.LGFunction) {
	r = map[string]map[string]lua.LGFunction{
		"context": map[string]lua.LGFunction{
			"set_cookie":       moduleHTTPContextSetCookie,
			"get_cookie":       moduleHTTPContextGetCookie,
			"set_charset":      moduleContexSetCharset,
			"set_header":       moduleContexSetHeader,
			"set_content_type": moduleHTTPContextSetContentType,
		},
		"url": map[string]lua.LGFunction{
			"encode": moduleURLEncode,
			"decode": moduleURLDecode,
		},
		"encoding": map[string]lua.LGFunction{
			"convert": moduleEncodingConvert,
		},
		"unicode": map[string]lua.LGFunction{
			"encode": moduleUnicodeEncode,
			"decode": moduleUnicodeDecode,
		},
		"md5": map[string]lua.LGFunction{
			"encrypt": moduleMd5Encrypt,
		},
		"des": map[string]lua.LGFunction{
			"encrypt":    moduleDESEncrypt,
			"decrypt":    moduleDESDecrypt,
			"qx_encrypt": moduleQXDESEncrypt,
			"qx_decrypt": moduleQXDESDecrypt,
		},
		"aes": map[string]lua.LGFunction{
			"encrypt": moduleAESEncrypt,
			"decrypt": moduleAESDecrypt,
		},
		"base64": map[string]lua.LGFunction{
			"encode":       moduleBase64Encode,
			"decode":       moduleBase64Decode,
			"encode_bytes": moduleBase64EncodeBytes,
			"decode_bytes": moduleBase64DecodeBytes,
		},
		"rsa": map[string]lua.LGFunction{
			"encrypt":   moduleRsaEncrypt,
			"decrypt":   moduleRsaDecrypt,
			"verify":    moduleRsaVerify,
			"make_sign": moduleRsaMakeSign,
		},
		"sha256": map[string]lua.LGFunction{
			"encrypt": moduleSha256Encrypt,
		},
		"sha1": map[string]lua.LGFunction{
			"encrypt": moduleSha1Encrypt,
		},
	}
	return
}
func getTypes() (r []*TypeBinder) {
	r = append(r, getHTTPClientTypeBinder())
	r = append(r, getMemcachedBinder())
	r = append(r, getWeixinTypeBinder())
	r = append(r, getinfluxTypeBinder())
	r = append(r, getMQTypeBinder())
	r = append(r, getImageDrawTypeBinder())
	return

}
