package parses

import (
	"strings"
)

type ITagMark interface {
	M3U8Type() string
}
type tagHandle func(tagData string, value string) ITagMark

var parseTags = map[string]tagHandle{
	"EXT-X-KEY": parseExtKey,
	"EXTINF":    parseExtInf,
}

type LineParse struct {
	Title string
	Data  string
}

func (l LineParse) Format() ITagMark {
	// fmt.Printf("%+v\n", l)
	for tagName, handle := range parseTags {
		prefix := "#" + tagName + ":"
		if strings.HasPrefix(l.Title, prefix) {
			tagData := l.Title[len(prefix):]
			return handle(tagData, l.Data)
		}
	}
	return nil
}
