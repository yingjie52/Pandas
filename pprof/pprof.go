package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// 主函数
func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	for {
		time.Sleep(1 * time.Second)
	}
}

//go tool pprof http://localhost:6060/debug/pprof/profile
//top
//list Eat
//web
//go tool pprof http://localhost:6060/debug/pprof/profile?seconds=60
//http://localhost:6060/debug/pprof/
