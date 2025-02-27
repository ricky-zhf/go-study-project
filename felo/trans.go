package felo

import (
	"StudyProject/common"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/afero"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"zlutils/consul"
)

//
//func init() {
//
//	consul.Init("172.30.43.72:8500", "me-srvconfs")
//}

const (
	downloadURL      = "https://test-admin.cocilive.com/api/excel/v1/download"
	saveFileName     = "felo/translate.xlsx"
	transOutDir      = "felo"
	transOutFilename = "translate.yaml"
)

const (
	MulLangServiceTag = "服务端Tag"
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

func SyncFeloMulLang(wg *sync.WaitGroup) {
	var (
		start = time.Now()
		err   error
	)
	defer func() {
		log.Printf("Sync felo multiple language to consul end. err:%v. cost:%v", err, time.Since(start))
		wg.Done()
	}()

	// 从google sheet从拉取多语
	data, err := parseExcel()
	if err != nil {
		log.Printf("Parse excel failed. err:%v", err)
		return
	}

	//if utils.IsEnvTest() {
	//	log.Printf("Get felo lang from google success. data=%+v", data)
	//}

	//同步到consul
	configPath := "me-srvconfs/common/felo-trans/"
	kv, _, err := consul.KV.Get(configPath, nil)
	if err != nil {
		log.Printf("Get data from consul failed. path:%v. err:%v", configPath, err)
		return
	}
	for key, value := range data {
		marshal, jsErr := json.Marshal(value)
		if jsErr != nil {
			log.Printf("Marshal data failed. err:%v. key:%v", jsErr, key)
			continue
		}

		//encoder, err := GzipJsonEncoder(marshal)
		//if err != nil {
		//	log.Printf("GzipJsonEncoder failed. err:%v", err)
		//	continue
		//}

		content := &api.KVPair{Key: kv.Key + "trans-" + key + ".json.gzip", Value: marshal}
		if _, err = consul.KV.Put(content, nil); err != nil {
			log.Printf("Consul put kv failed. key:%v. path:%v. err:%v", key, configPath, err)
			return
		}
	}
	return
}

type MulLanguages map[string]string

type LanguageData map[string]MulLanguages

type ServiceMulLang map[string]LanguageData

func parseExcel() (serviceMulLang ServiceMulLang, err error) {
	fileUrl := "https://docs.google.com/spreadsheets/d/1StmDHmnNbB9DofeIGDSwp9GxJM2KyxvtiQyVdr9PPao/export?format=xlsx"
	resp, err := http.Get(fileUrl)
	if err != nil {
		log.Printf("Get felo multiple language failed. err:%v", err)
		return
	}
	defer resp.Body.Close()

	f, _ := excelize.OpenReader(resp.Body)

	serviceMulLang = make(ServiceMulLang)
	for _, sheetName := range f.GetSheetMap() {
		rows, _ := f.GetRows(sheetName)
		if rows == nil || len(rows) == 0 {
			continue
		}

		headers := rows[0]
		serverTagCol := -1
		for i, header := range headers {
			if header == common.MulLangServiceTag {
				serverTagCol = i
				break
			}
		}
		if serverTagCol == -1 {
			log.Printf("sheet doesn't have service tag. sheet:%v", sheetName)
			continue
		}

		for _, row := range rows[1:] {
			if len(row) <= serverTagCol || row[serverTagCol] == "" || row[serverTagCol] == common.MulLangServiceTag {
				continue
			}

			parseRows(row, headers, serverTagCol, serviceMulLang)
		}
	}
	return
}

func parseRows(row []string, headers []string, serverTagCol int, servers ServiceMulLang) {
	item := make(map[string]string)
	for colIdx, cell := range row {
		header := headers[colIdx]
		item[header] = cell
	}

	serviceNames := strings.Split(row[serverTagCol], ",")
	for _, serviceName := range serviceNames {
		serviceName = strings.TrimSpace(serviceName)
		server, exists := servers[serviceName]
		if !exists {
			server = make(LanguageData)
		}

		key := item["Key"]
		for langKey, langHeader := range common.FeloLanguageHeaders {
			langValue, has := item[langHeader]
			if !has || langValue == "" {
				continue
			}

			if _, has := server[key]; !has {
				server[key] = make(MulLanguages)
			}

			server[key][langKey] = langValue
		}

		// 更新服务器信息
		servers[serviceName] = server
	}
}

// 写入 YAML 文件
// writeYAML 将 serviceItems 转换为 YAML 文件并写入
//func writeYAML(serviceItems map[string]ServiceMulLang) error {
//	for serviceName, service := range serviceItems {
//		serviceDir := filepath.Join(transOutDir, serviceName)
//		if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
//			if err := os.MkdirAll(serviceDir, 0755); err != nil {
//				return err
//			}
//		}
//
//		transFile := filepath.Join(serviceDir, transOutFilename)
//		transFileItems := make(map[string]map[string]string)
//
//		if _, err := os.Stat(transFile); err == nil {
//			data, err := ioutil.ReadFile(transFile)
//			if err != nil {
//				return err
//			}
//			yaml.Unmarshal(data, &transFileItems)
//		}
//
//		// 将每个语言和区域数据添加到 transFileItems
//		for langKey, lang := range service.Languages {
//			for areaKey, area := range lang.Areas {
//				fileItem, exists := transFileItems[areaKey]
//				if !exists {
//					fileItem = make(map[string]string)
//					fileItem["en"] = area.Key
//				}
//				fileItem[langKey] = area.Value
//				transFileItems[areaKey] = fileItem
//			}
//		}
//
//		// 写入 YAML 文件
//		data, err := yaml.Marshal(transFileItems)
//		if err != nil {
//			return err
//		}
//		if err := ioutil.WriteFile(transFile, data, 0644); err != nil {
//			return err
//		}
//	}
//	return nil
//}

func run() {
	// 下载 Excel 文件
	//if err := downloadFile(downloadURL, saveFileName); err != nil {
	//	log.Fatalf("Download error: %v", err)
	//}

	// 解析 Excel 文件
	//serviceItems, err := parseExcel(saveFileName)
	//if err != nil {
	//	log.Fatalf("Parse error: %v", err)
	//}

	// 写入 YAML 文件
	//if err = writeYAML(serviceItems); err != nil {
	//	log.Fatalf("YAML write error: %v", err)
	//}

	// 从google sheet从拉取多语
	var wg sync.WaitGroup
	wg.Add(1)
	SyncFeloMulLang(&wg)
	wg.Wait()

	return
}

func GzipJsonEncoder(in []byte) ([]byte, error) {
	buf := bytes.Buffer{}
	gw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	gw.Write(in)
	defer gw.Close()
	return buf.Bytes(), nil
}
