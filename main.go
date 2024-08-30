package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	// 初始化cron调度器
	c := cron.New(cron.WithSeconds())

	// 定义一个任务时间，5秒后执行
	taskTime := time.Now().Add(5 * time.Second).Format("05 04 15 02 01 *")

	// 声明一个变量来存储entryID
	var entryID cron.EntryID

	// 添加任务并获取entryID
	entryID, err := c.AddFunc(taskTime, func() {
		// 任务被执行时打印entryID
		fmt.Println("任务执行, EntryID:", entryID)
	})
	if err != nil {
		fmt.Println("添加任务失败:", err)
		return
	}

	// 打印已获取到的entryID
	fmt.Println("已添加任务, EntryID:", entryID)

	// 启动调度器
	c.Start()

	// 阻塞主线程，等待任务执行
	select {}
}
