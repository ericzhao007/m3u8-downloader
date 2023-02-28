package m3u8

import (
	"bufio"
	"github.com/ericzhao007/m3u8-downloader/m3u8/parses"
	"io"
	"strings"
)

type M3U8Parse struct {
	lines      []parses.ITagMark
	Encryption *parses.M3U8Encryption
	PayList    []*parses.M3U8Inf
}

func NewM3U8Parse(data io.Reader) *M3U8Parse {
	sc := bufio.NewScanner(data)
	multiFlag := false
	mp := &M3U8Parse{}
	var (
		lineParse parses.LineParse
	)
	for sc.Scan() {
		line := sc.Text()
		if len(line) == 0 {
			continue
		}
		if line[:1] == "#" {
			lineParse.Title = line
			//TODO 只有指定的标签才是多行
			if strings.HasPrefix(line, "#EXTINF") {
				multiFlag = true
				continue
			} else {
				multiFlag = false
				if parseRes := lineParse.Format(); parseRes != nil {
					mp.lines = append(mp.lines, parseRes)
				}
				lineParse.Title = ""
				lineParse.Data = ""
				continue
			}

		}
		if multiFlag {
			lineParse.Data = line
			if parseRes := lineParse.Format(); parseRes != nil {
				mp.lines = append(mp.lines, parseRes)
			}
			lineParse.Title = ""
			lineParse.Data = ""
			multiFlag = false
		}
	}
	return mp
}

func (mp *M3U8Parse) Do() {
	for _, line := range mp.lines {
		m3u8Type := line.M3U8Type()
		switch m3u8Type {
		case "m3u8encryption":
			mp.Encryption = line.(*parses.M3U8Encryption)
		case "m3u8inf":
			if inf, ok := line.(*parses.M3U8Inf); ok {
				mp.PayList = append(mp.PayList, inf)
			}
		}
	}

}

func ParseM3U8(data io.Reader) *M3U8Parse {
	mp := NewM3U8Parse(data)
	mp.Do()
	return mp
}
