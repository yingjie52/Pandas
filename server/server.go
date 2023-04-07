package main

import (
	"encoding/json"
	"net/http"
)

//import (
//	"fmt"
//
//	"github.com/valyala/fasthttp"
//)
//
//// 用fasthttp实现的http服务器
//func main() {
//	handler := func(ctx *fasthttp.RequestCtx) {
//		fmt.Fprintf(ctx, "Hello, world! Requested path is %q", ctx.Path())
//	}
//
//	if err := fasthttp.ListenAndServe(":8081", handler); err != nil {
//		fmt.Println(err)
//	}
//}

//func main() {
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprintf(w, "Hello, world! Requested path is %q", r.URL.Path)
//	})
//
//	err := http.ListenAndServe(":8081", nil)
//	if err != nil {
//		fmt.Println("Error starting server:", err.Error())
//	}
//}

type HeightResponse struct {
	Height string `json:"height"`
}

func main() {
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		name := r.URL.Query().Get("name")

		// 查询身高信息
		var height string
		if name == "xiaming" {
			height = "172cm"
		} else {
			height = "unknown"
		}

		// 返回响应数据
		resp := HeightResponse{Height: height}
		json.NewEncoder(w).Encode(resp)
	})

	http.ListenAndServe(":8080", nil)
}
