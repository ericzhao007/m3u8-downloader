package parses

import (
	"strings"

	"github.com/duke-git/lancet/v2/cryptor"
)

type M3U8Encryption struct {
	Method     string
	IV         string
	URI        string
	key        []byte
	originData map[string]string
}

func parseExtKey(tagData string, value string) ITagMark {
	mk := &M3U8Encryption{
		originData: make(map[string]string),
	}
	kvs := strings.Split(tagData, ",")
	for _, kv := range kvs {
		kvres := strings.Split(kv, "=")
		if len(kvres) > 1 {
			mk.originData[kvres[0]] = strings.Trim(kvres[1], "\"")
		}
	}
	if method, exists := mk.originData["METHOD"]; exists {
		mk.Method = method
	}
	if iv, exists := mk.originData["IV"]; exists {
		mk.IV = iv
	}
	if uri, exists := mk.originData["URI"]; exists {
		mk.URI = uri
	}
	return mk
}

func (mk *M3U8Encryption) M3U8Type() string {
	return "m3u8encryption"
}

func (mk *M3U8Encryption) SetKey(key []byte) {
	mk.key = key
}

func (mk *M3U8Encryption) GetKey() []byte {
	return mk.key
}

func (mk *M3U8Encryption) AesDecrypt(data []byte) []byte {
	return cryptor.AesCbcDecrypt(data, mk.key)
}
