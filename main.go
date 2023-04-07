package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"net/http"
)

func main() {
	// 打开 Excel 文件
	f, err := excelize.OpenFile("123.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取第一个 sheet 并读取数据
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/query", nil)
	if err != nil {
		fmt.Print(12345)
	}
	query := req.URL.Query()

	// 遍历每一行数据
	for i, row := range rows {
		if i == 0 { // 跳过表头
			continue
		}
		name123 := row[4]

		result2 := gjson.Parse(name123)

		//	获取result2的的key和
		result2.ForEach(func(key, value gjson.Result) bool {
			query.Add(key.String(), value.String())
			return true

		})

	}
	req.URL.RawQuery = query.Encode()

	// 发起 HTTP 请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("123")

	}
	defer resp.Body.Close()
	//	// 3. 处理响应
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

}

//func main() {
//	req, err := http.NewRequest("GET", "http://localhost:8080/query", nil)
//	if err != nil {
//		fmt.Print(12345)
//	}
//
//	query := req.URL.Query()
//	query.Add("name", "xiaming")
//	req.URL.RawQuery = query.Encode()
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		fmt.Print(1234567)
//	}
//
//	var heightResp HeightResponse
//	err = json.NewDecoder(resp.Body).Decode(&heightResp)
//	if err != nil {
//		fmt.Print(12345678)
//	}
//
//	fmt.Print(heightResp.Height)
//}

git remote add origin https://gitee.com/huangyingjie6/pandas.git

