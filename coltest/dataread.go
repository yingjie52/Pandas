package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"strings"
)

// 遍历当前文件下的excel文件,并返回文件名
func readExcelFile() []string {
	var data []string
	//打开文件夹
	dir, err := os.Open("./")
	if err != nil {
		fmt.Println(err)
	}
	//读取文件夹下的文件
	fis, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
	}
	//遍历文件夹下的文件
	for _, fi := range fis {
		//判断是否为文件夹
		if fi.IsDir() {
			continue
		}
		//判断是否为excel文件
		if strings.HasSuffix(fi.Name(), ".xlsx") {
			data = append(data, fi.Name())
		}
	}
	return data
}

var data [][]interface{}

// 读取Excel文件中的数据到内存中
func ReadExcelFileToMemory() ([][]interface{}, error) {

	// 获取excel文件名称
	filename := readExcelFile()
	// 打开Excel文件
	f, err := excelize.OpenFile(filename[0])
	if err != nil {
		return nil, err
	}

	//遍历所有的sheet
	for _, sheet := range f.GetSheetMap() {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, err
		}
		// 遍历每一行
		for _, row := range rows {
			var rowData []interface{}

			// 获取每一行中的所有单元格
			for _, cell := range row {
				// 将单元格中的数据存储到缓存中
				rowData = append(rowData, cell)
			}
			// 将一行数据存储到内存中
			data = append(data, rowData)
		}
	}

	return data, nil
}

// 清空内存保存的数据
func ClearMemory(data *[][]interface{}) {
	*data = make([][]interface{}, 0)
	//clearMemory(&data)
}

var datach = make(chan *request)

// 将数据写入结构体中，然后将结构体写入通道中
func dataload() {

	data, _ := ReadExcelFileToMemory()
	con, _ := data[1][1].(int)        // 用例编号
	methods, _ := data[1][2].(string) // 请求方法
	urls, _ := data[1][3].(string)
	params, _ := data[1][4].(string) // 请求参数

	r := &request{
		cons:   con,
		method: methods,
		url:    urls,
		param:  params,
	}

	for i := 0; i < 5; i++ {
		datach <- r
	}

}
