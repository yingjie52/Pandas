package main

import (
	"fmt"
	"sync"
	"time"
)

// header 打印表头信息
func header() {
	// 打印的时长都为毫秒 总请数
	fmt.Println(" 耗时  并发数  成功数  失败数  qps  最长耗时  最短耗时  平均耗时")
	return
}

// table 打印表格
func table(successNum, failureNum uint64, qps, averageTime, maxTimeFloat, minTimeFloat, requestTimeFloat float64, chanIDLen int) {
	// 打印的时长都为毫秒
	result := fmt.Sprintf("%4.0fs│%7d│%7d│%7d│%8.2f│%8.2f│%8.2f│%8.2f│%8s│%8s│%v",
		requestTimeFloat, chanIDLen, successNum, failureNum, qps, maxTimeFloat, minTimeFloat, averageTime)
	fmt.Println(result)
	return
}

// calculateData 计算数据
func calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum uint64,
	chanIDLen int) {
	if processingTime == 0 {
		processingTime = 1
	}
	var (
		qps              float64
		averageTime      float64
		maxTimeFloat     float64
		minTimeFloat     float64
		requestTimeFloat float64
	)
	// 平均 QPS 成功数*总协程数/总耗时 (每秒)
	if processingTime != 0 {
		qps = float64(successNum*1e9*concurrent) / float64(processingTime)
	}
	// 平均时长 总耗时/总请求数/并发数 纳秒=>毫秒
	if successNum != 0 && concurrent != 0 {
		averageTime = float64(processingTime) / float64(successNum*1e6)
	}
	// 纳秒=>毫秒
	maxTimeFloat = float64(maxTime) / 1e6
	minTimeFloat = float64(minTime) / 1e6
	requestTimeFloat = float64(requestTime) / 1e9
	// 打印的时长都为毫秒
	table(successNum, failureNum, qps, averageTime, maxTimeFloat, minTimeFloat, requestTimeFloat, chanIDLen)
}

type RequestResults struct {
	ID            string // 消息ID
	ChanID        uint64 // 消息ID
	Time          uint64 // 请求时间 纳秒
	IsSucceed     bool   // 是否请求成功
	ErrCode       int    // 错误码
	ReceivedBytes int64
}

// ReceivingResults 接收结果并处理
// 统计的时间都是纳秒，显示的时间 都是毫秒
// concurrent 并发数

var requestTimeList []uint64

func ReceivingResults(concurrent uint64, ch <-chan *RequestResults, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	var stopChan = make(chan bool)
	// 时间
	var (
		processingTime uint64 // 处理总时间
		requestTime    uint64 // 请求总时间
		maxTime        uint64 // 最大时长
		minTime        uint64 // 最小时长
		successNum     uint64 // 成功处理数，code为0
		failureNum     uint64 // 处理失败数，code不为0
		chanIDLen      int    // 并发数
		mutex          = sync.RWMutex{}
	)
	statTime := uint64(time.Now().UnixNano())
	// 错误码/错误个数

	// 定时输出一次计算结果
	ticker := time.NewTicker(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				endTime := uint64(time.Now().UnixNano())
				mutex.Lock()
				go calculateData(concurrent, processingTime, endTime-statTime, maxTime, minTime, successNum, failureNum, chanIDLen)
				mutex.Unlock()
			case <-stopChan:
				// 处理完成
				return
			}
		}
	}()
	header()
	for data := range ch {
		mutex.Lock()
		// fmt.Println("处理一条数据", data.ID, data.Time, data.IsSucceed, data.ErrCode)
		processingTime = processingTime + data.Time
		if maxTime <= data.Time {
			maxTime = data.Time
		}
		if minTime == 0 {
			minTime = data.Time
		} else if minTime > data.Time {
			minTime = data.Time
		}
		// 是否请求成功
		if data.IsSucceed == true {
			successNum = successNum + 1
		} else {
			failureNum = failureNum + 1
		}
		requestTimeList = append(requestTimeList, data.Time)
		mutex.Unlock()
	}
	// 数据全部接受完成，停止定时输出统计数据
	stopChan <- true
	endTime := uint64(time.Now().UnixNano())
	requestTime = endTime - statTime
	calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, chanIDLen)

	fmt.Printf("\n\n")
	fmt.Println("*************************  结果 stat  ****************************")
	fmt.Println("处理协程数量:", concurrent)
	// fmt.Println("处理协程数量:", concurrent, "程序处理总时长:", fmt.Sprintf("%.3f", float64(processingTime/concurrent)/1e9), "秒")
	fmt.Println("请求总数（并发数*请求数 -c * -n）:", successNum+failureNum, "总请求时间:",
		fmt.Sprintf("%.3f", float64(requestTime)/1e9),
		"秒", "successNum:", successNum, "failureNum:", failureNum)
	fmt.Println("*************************  结果 end   ****************************")
	fmt.Printf("\n\n")
}
