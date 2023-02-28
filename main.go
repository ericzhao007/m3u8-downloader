package main

import (
	"flag"
	"fmt"

	"github.com/ericzhao007/m3u8-downloader/spiders"
)

var m3u8Url = flag.String("m3u8", "", "m3u8下载链接")
var useMemory = flag.Bool("mem", false, "是否内存加速，会占用大量内存")
var workerNum = flag.Int("m", 0, "并发数")
var filePath = flag.String("f", "movie.mp4", "保存的文件地址")

func main() {
	flag.Parse()
	if *m3u8Url == "" {
		fmt.Println("请填写m3u8地址")
		return
	}
	// appLog, err := os.Create("app.log")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer appLog.Close()
	// log.Default().SetOutput(appLog)
	var withs []spiders.WithFunc
	if *useMemory {
		withs = append(withs, spiders.WithMemoryStoreEngine())
	} else {
		withs = append(withs, spiders.WithDiskStoreEngine())
	}
	withs = append(withs, spiders.WithWorkerNum(*workerNum))
	s := spiders.NewSpiders(*m3u8Url, withs...)
	if err := s.Run(*filePath); err != nil {
		fmt.Println(err)
	}

}
