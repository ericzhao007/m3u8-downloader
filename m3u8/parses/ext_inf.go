package parses

import (
	"github.com/duke-git/lancet/v2/cryptor"
	"strings"
)

type M3U8Inf struct {
	Duration string
	Title    string
	Data     string
}

func parseExtInf(tagData string, value string) ITagMark {
	mi := &M3U8Inf{}
	kvs := strings.Split(tagData, ",")
	if len(kvs) > 0 {
		mi.Duration = kvs[0]
	}
	if len(kvs) > 1 {
		mi.Title = kvs[1]
	}
	mi.Data = value
	return mi
}

func (mi *M3U8Inf) M3U8Type() string {
	return "m3u8inf"
}

func (mi *M3U8Inf) Name() string {
	return mi.Data
}

func (mi *M3U8Inf) HashId() string {
	return cryptor.Md5String(mi.Data)
}
