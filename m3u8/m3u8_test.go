package m3u8

import (
	"bytes"
	"github.com/ericzhao007/m3u8-downloader/m3u8/parses"
	"io"
	"testing"
)

func TestM3U8Parse_Do(t *testing.T) {
	type fields struct {
		lines      []parses.ITagMark
		Encryption *parses.M3U8Encryption
		PayList    []*parses.M3U8Inf
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &M3U8Parse{
				lines:      tt.fields.lines,
				Encryption: tt.fields.Encryption,
				PayList:    tt.fields.PayList,
			}
			mp.Do()
		})
	}
}

func TestNewM3U8Parse(t *testing.T) {
	type args struct {
		data io.Reader
	}
	tests := []struct {
		name string
		args args
		want *M3U8Parse
	}{
		{
			name: "base test",
			args: args{
				data: bytes.NewReader([]byte(`
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-ALLOW-CACHE:YES
#EXT-X-TARGETDURATION:5
#EXT-X-MEDIA-SEQUENCE:0
#EXT-X-PLAYLIST-TYPE:VOD
#EXT-X-KEY:METHOD=AES-128,URI="key.key",IV=0x00000000000000000000000000000000
#EXTINF:4.840000,
N2Y0ZTFmMj0.ts
#EXTINF:5.120000,
N2Y0ZTFmMj1.ts
#EXTINF:5.000000,
N2Y0ZTFmMj2.ts
#EXTINF:5.000000,
N2Y0ZTFmMj3.ts
#EXTINF:5.000000,
N2Y0ZTFmMj4.ts
#EXTINF:5.000000,
N2Y0ZTFmMj5.ts
#EXTINF:5.000000,
N2Y0ZTFmMj6.ts
#EXTINF:4.960000,
N2Y0ZTFmMj7.ts
#EXTINF:5.000000,
				`)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewM3U8Parse(tt.args.data)
			got.Do()
			t.Logf("%+v", got)
			// if got := NewM3U8Parse(tt.args.data); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("NewM3U8Parse() = %v, want %v", got, tt.want)
			// }
		})
	}
}
