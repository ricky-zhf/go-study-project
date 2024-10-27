package felo

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	downloadURL      = "https://test-admin.cocilive.com/api/excel/v1/download"
	saveFileName     = "felo/translate.xlsx"
	transOutDir      = "felo"
	transOutFilename = "translate.yaml"
)

var languageHeaders = map[string]string{
	"en": "English(en)",
	"ar": "Arabic(ar)",
	"tr": "Turkish(tr)",
	"zh": "Chinese(zh-CN)",
}

var languageOrder = []string{"zh", "en", "ar", "tr"}

// 下载文件
func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dir := filepath.Dir(dest)
	fs := afero.NewOsFs()
	if err = fs.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 写入文件
	err = afero.WriteFile(fs, dest, body, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// 解析 Excel 文件
func parseExcel(filePath string) (map[string]map[string]map[string]string, error) {
	filePath = "/Users/zhouhuaifeng/GoWorkspace/src/StudyProject/felo/felo/translate.xlsx"
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// service_name - language_key - k - v
	serviceItems := make(map[string]map[string]map[string]string)

	for _, sheetName := range f.GetSheetList() {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return nil, err
		}

		var (
			headers      = rows[0]
			serverTagCol = -1
		)

		for i, header := range headers {
			if header != "服务端Tag" {
				continue
			}
			serverTagCol = i
			break
		}
		if serverTagCol == -1 {
			// add log
			return nil, errors.New("Sheet doesn't have server tag.")
		}

		for _, row := range rows[1:] {
			if len(row) <= serverTagCol || len(row[serverTagCol]) == 0 || row[serverTagCol] == "服务端Tag" {
				continue
			}

			// 服务端翻译
			var (
				item        = make(map[string]string)
				serverNames = row[serverTagCol]
			)

			for colIdx, cell := range row {
				header := headers[colIdx]
				item[header] = cell
			}

			for _, serviceName := range strings.Split(serverNames, ",") {
				serviceName = strings.TrimSpace(serviceName)
				if _, exists := serviceItems[serviceName]; !exists {
					serviceItems[serviceName] = make(map[string]map[string]string)
				}
				serviceItems[serviceName][item["Key"]] = item
			}
		}
	}
	return serviceItems, nil
}

// 写入 YAML 文件
func writeYAML(serviceItems map[string]map[string]map[string]string) error {
	for serviceName, items := range serviceItems {
		serviceDir := filepath.Join(transOutDir, serviceName)
		if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
			if err := os.MkdirAll(serviceDir, 0755); err != nil {
				return err
			}
		}

		transFile := filepath.Join(serviceDir, transOutFilename)
		transFileItems := make(map[string]map[string]string)

		if _, err := os.Stat(transFile); err == nil {
			data, err := ioutil.ReadFile(transFile)
			if err != nil {
				return err
			}
			yaml.Unmarshal(data, &transFileItems)
		}

		for key, item := range items {
			fileItem := make(map[string]string)
			fileItem["en"] = key
			for _, lang := range languageOrder {
				header := languageHeaders[lang]
				if langValue, exists := item[header]; exists && langValue != "" {
					fileItem[lang] = langValue
				}
			}
			transFileItems[key] = fileItem
		}

		data, err := yaml.Marshal(transFileItems)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(transFile, data, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func run() {
	// 下载 Excel 文件
	//if err := downloadFile(downloadURL, saveFileName); err != nil {
	//	log.Fatalf("Download error: %v", err)
	//}

	// 解析 Excel 文件
	serviceItems, err := parseExcel(saveFileName)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	// 写入 YAML 文件
	if err = writeYAML(serviceItems); err != nil {
		log.Fatalf("YAML write error: %v", err)
	}

	fmt.Println("Translation files generated successfully.")
}
