package schedule

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// 1.设置任务列表、执行者数量
// 2.运行

type ScheduleTask struct {
	Name      string
	workerNum int
	taskList  []any
}

func NewScheduleTask(name string, taskList []any, workerNum int) *ScheduleTask {
	if workerNum <= 0 {
		workerNum = runtime.NumCPU()
	}
	st := &ScheduleTask{
		Name:      name,
		workerNum: workerNum,
		taskList:  taskList,
	}
	return st
}

func (st *ScheduleTask) Run(callFunc func(v any, bar *mpb.Bar)) {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	p := mpb.NewWithContext(ctx, mpb.WithWaitGroup(&wg))
	//p := mpb.New(mpb.WithWaitGroup(&wg))
	// wg.Add(1)
	totalBar := p.AddBar(int64(len(st.taskList)),
		mpb.PrependDecorators(
			decor.Name("downloading", decor.WC{W: 12, C: decor.DidentRight}),
			decor.Name(st.Name, decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.Percentage(decor.WC{W: 5}), "done"),
		),
	)
	taskChan := make(chan any, st.workerNum)
	wg.Add(st.workerNum)
	for i := 0; i < st.workerNum; i++ {
		go func(i int) {
			defer wg.Done()
			taskName := fmt.Sprintf("Task#%02d:", i)
			bar := p.AddBar(0,
				mpb.BarID(i),
				mpb.PrependDecorators(
					decor.Any(func(s decor.Statistics) string {
						return taskName
					}, decor.WCSyncSpaceR),
					decor.CountersKibiByte("% .2f / % .2f"),
				),
				mpb.AppendDecorators(
					decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
					decor.OnComplete(decor.Percentage(decor.WC{W: 5}), "done"),
				),
			)
			for {
				v, ok := <-taskChan
				if !ok {
					break
				}
				tname := ""
				if vname, ok := v.(interface {
					Name() string
				}); ok {
					tname = vname.Name()
				}
				taskName = fmt.Sprintf("Task#%02d:%s", i, tname)
				callFunc(v, bar)
				bar.SetCurrent(0)
				totalBar.Increment()
			}

		}(i)
	}
	go func() {
		for _, task := range st.taskList {
			taskChan <- task
		}
		close(taskChan)
	}()
	totalBar.Wait()
	cancel()

	p.Wait()
}
