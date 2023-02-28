package spiders

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ericzhao007/m3u8-downloader/m3u8"
	"github.com/ericzhao007/m3u8-downloader/m3u8/parses"
	"github.com/ericzhao007/m3u8-downloader/schedule"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/vbauerster/mpb/v8"
)

type IStore interface {
	Load(key string) (value []byte, exists bool)
	Store(key string, value []byte) error
	Clear() error
}

type WithFunc func(s *Spiders)

type Spiders struct {
	m3u8Url   *url.URL
	tsData    IStore
	workerNum int
}

func NewSpiders(m3u8Url string, withs ...WithFunc) *Spiders {
	baseUrl, err := url.Parse(m3u8Url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	s := &Spiders{
		m3u8Url: baseUrl,
	}
	for _, withFunc := range withs {
		withFunc(s)
	}
	if s.tsData == nil {
		s.tsData = NewMemStoreEngine()
	}
	return s
}

func (s *Spiders) parseM3u8(m3u8Url string) (*m3u8.M3U8Parse, error) {
	resp, err := http.DefaultClient.Get(m3u8Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	mp := m3u8.ParseM3U8(resp.Body)
	if mp.Encryption != nil {
		u, err := url.Parse(mp.Encryption.URI)
		if err != nil {
			return nil, err
		}
		encryKeyResp, err := http.DefaultClient.Get(s.m3u8Url.ResolveReference(u).String())
		if err != nil {
			return nil, err
		}
		defer encryKeyResp.Body.Close()
		aesKey, err := ioutil.ReadAll(encryKeyResp.Body)
		if err != nil {
			return nil, err
		}
		mp.Encryption.SetKey(aesKey)

	}
	return mp, nil
}

func (s *Spiders) Run(filePath string) error {
	// 解析入口文件
	mp, err := s.parseM3u8(s.m3u8Url.String())
	if err != nil {
		return err
	}
	taskList := make([]any, 0, len(mp.PayList))
	for _, payItem := range mp.PayList {
		taskList = append(taskList, payItem)
	}
	// 调度器
	st := schedule.NewScheduleTask(filePath, taskList, s.workerNum)
	st.Run(func(v any, bar *mpb.Bar) {
		m3u8inf := v.(*parses.M3U8Inf)
		m3u8infUrl, _ := url.Parse(m3u8inf.Data)
		dataUrl := s.m3u8Url.ResolveReference(m3u8infUrl)
		resp, _ := http.DefaultClient.Get(dataUrl.String())
		defer resp.Body.Close()
		bar.SetTotal(resp.ContentLength, false)
		proxyRader := bar.ProxyReader(resp.Body)
		defer proxyRader.Close()
		bys, _ := ioutil.ReadAll(proxyRader)
		if mp.Encryption != nil {
			bys = mp.Encryption.AesDecrypt(bys)
		}
		s.tsData.Store(m3u8inf.Name(), bys)
	})
	// 视频合并
	videoFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer videoFile.Close()
	fmt.Println("正在合并...")
	for _, payItem := range mp.PayList {
		v, exists := s.tsData.Load(payItem.Name())
		if !exists {
			return errors.New("download faild :" + payItem.Name())
		}
		_, err := io.Copy(videoFile, bytes.NewReader(v))
		if err != nil {
			return fmt.Errorf("copy faild, %s %w", payItem.Name(), err)
		}
	}
	defer fmt.Println("下载完成：", filePath)
	return s.tsData.Clear()
}
