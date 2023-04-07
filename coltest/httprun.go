package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type request struct {
	cons   int
	method string
	url    string
	param  string
}

//func Httprun(chan2 chan *request) {
//
//	var wg sync.WaitGroup
//	rd := <-chan2
//	pool, _ := ants.NewPool(6)
//
//	starttime := time.Now()
//
//	for i := 0; i < 99999999999; i++ {
//		wg.Add(1)
//		err := pool.Submit(func() {
//			defer wg.Done()
//
//			// 1. 创建请求
//			req, err := http.NewRequest(rd.method, rd.url, nil)
//			if err != nil {
//				log.Fatal(err)
//			}
//			// 2. 发送请求
//			resp, err := http.DefaultClient.Do(req)
//			if err != nil {
//				log.Fatal(err)
//			}
//			// 3. 处理响应
//			defer func(Body io.ReadCloser) {
//				err := Body.Close()
//				if err != nil {
//
//				}
//			}(resp.Body)
//			body, err := io.ReadAll(resp.Body)
//			if err != nil {
//				log.Fatal(err)
//			}
//			fmt.Print(string(body), "\n")
//		})
//		if err != nil {
//			return
//		}
//		time.Sleep(1 * time.Second)
//		if time.Since(starttime) > 5*time.Second {
//			return
//		}
//
//	}
//	wg.Wait()
//	pool.Release()
//}

func Httprun(chan2 chan *request) {

	var wg sync.WaitGroup

	pool, _ := ants.NewPool(6)

	stopchan := make(chan bool)

	go func() {
		starttime := time.Now()
		for {
			rd := <-chan2
			wg.Add(1)
			err := pool.Submit(func() {
				defer wg.Done()

				// 1. 创建请求
				req, err := http.NewRequest(rd.method, rd.url, nil)
				if err != nil {
					log.Fatal(err)
				}

				result2 := gjson.Parse(rd.param)

				query := req.URL.Query()

				//	获取result2的的key和value
				result2.ForEach(func(key, value gjson.Result) bool {
					query.Add(key.String(), value.String())
					return true

				})

				req.URL.RawQuery = query.Encode()

				// 2. 发送请求
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				// 3. 处理响应
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {

					}
				}(resp.Body)

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Print(string(body), "\n")
			})
			if err != nil {
				return
			}
			//停止10毫秒
			time.Sleep(500 * time.Millisecond)
			if time.Since(starttime) > 5*time.Second {
				//运行时长
				fmt.Println(time.Since(starttime))
				stopchan <- true
				return
			}
		}

	}()
	<-stopchan
	wg.Wait()
	pool.Release()
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	go dataload()
	Httprun(datach)
}
