package main

import (
	"fmt"
	"math/rand"
	"time"
)

//数据保存
var LimitQueue map[string][]int64

func LimitFreqSingle(queueName string, count uint, timeWindow int64) bool {
	currTime := time.Now().Unix()

	//创建
	if _, ok := LimitQueue[queueName]; !ok {
		LimitQueue[queueName] = make([]int64, 0)
	}

	//队列未满
	if uint(len(LimitQueue[queueName])) < count {
		LimitQueue[queueName] = append(LimitQueue[queueName], currTime)
		return true
	}

	//队列满了
	//找到最早的那个时间
	firstTime := LimitQueue[queueName][0]
	//小于窗口的限定时间，说明还没过期，此次请求不允许通过
	if currTime-firstTime <= timeWindow {
		return false
	}

	//校验通过，将本次请求添加到队列中
	LimitQueue[queueName] = LimitQueue[queueName][1:]
	LimitQueue[queueName] = append(LimitQueue[queueName], currTime)

	return true
}

func init() {
	LimitQueue = make(map[string][]int64)
}

func main() {
	//timeWindow秒内，可以容纳count次请求

	req := "index"
	var timeWindow int64 = 10
	var count uint = 5
	runSec := 0

	for i := 0; i < 100; i++ {

		res := LimitFreqSingle(req, count, timeWindow)

		fmt.Printf("第%d次请求，时间：%d，是否通行：%t\n", i, runSec, res)
		sec := rand.Intn(5)
		runSec += sec
		time.Sleep(time.Duration(sec) * time.Second)
	}
}
