package felo

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
)

func tkdo() {
	// 打开 Excel 文件
	f, err := excelize.OpenFile("/Users/zhouhuaifeng/Downloads/higo埋点盘点 (1).xlsx")
	if err != nil {
		log.Fatalf("Failed to open Excel file: %v", err)
	}
	defer f.Close()

	// 存放结果的 map
	result := make(map[string]map[string]string)

	// 读取 "felo需要保留埋点" sheet
	rows, err := f.GetRows("felo需要保留埋点")
	if err != nil {
		log.Fatalf("Failed to get rows: %v", err)
	}

	// 遍历行，假设第1列为事件名，第2列为page_name
	var eventName string = "todo"
	for i, row := range rows[1:] { // 跳过表头
		if len(row) < 5 {
			continue // 跳过不完整的数据行
		}
		if len(row[3]) != 0 {
			eventName = row[3] // 事件名
		}

		pageName := row[4] // page_name
		if len(pageName) == 0 || pageName == "-" {
			pageName = "todo===" + fmt.Sprintf("%v", i)
		}
		// 创建一个 map 存储 event
		result[pageName] = map[string]string{"event": eventName}
	}

	// 转换为 JSON 并写入文件
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// 写入到文件
	if err = os.WriteFile("output.json", jsonData, 0644); err != nil {
		log.Fatalf("Failed to write JSON file: %v", err)
	}

	fmt.Println("JSON data written to output.json successfully.")
}
