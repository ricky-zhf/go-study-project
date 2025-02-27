package saas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	jp   = "https://xz-admin-jp.xiaozlife.com"
	na   = "https://xz-admin-na.wjacloud.com"
	sg   = "https://xz-admin-sg.wjacloud.com"
	eu   = "https://xz-admin-eu.xiaozlife.com"
	path = "/admin/v1/resource/updateLivePage"
	//"/admin/v1/resource/updateLivePage"
	//urls = []string{jp + path, na + path, sg + path, eu + path}
	urls = []string{sg + path}

	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiMTk1ZDE3ZWEtYmFjMi00OWE5LTk2NzItNTdjMjZiMmMxNmRiIiwiSUQiOjEsIlVzZXJuYW1lIjoiYWRtaW4iLCJOaWNrTmFtZSI6IueuoeeQhuWRmCIsIkF1dGhvcml0eUlkIjo4ODgsIkJ1ZmZlclRpbWUiOjg2NDAwLCJleHAiOjE3MjE5NjAzMjgsImlzcyI6InFtUGx1cyIsIm5iZiI6MTcyMTM1NDUyOH0.KT92kDH-yuUm7HggbVBNCnEB94cB-q6nLAfG5Xd9UpQ"
	str   = "UPDATE xz_server_manager.category SET app_id=0 WHERE id IN (21,22,25,26,27);"
)

// 假设这是从某处获取的SQL更新语句数组
//var sqlUpdates = []string{}

func run() {
	// 将SQL语句转换为JSON格式
	sql := strings.Split(str, "\n")
	//fmt.Println(sql)
	fmt.Println("sql len=", len(sql))

	for _, url := range urls {
		fmt.Println("start req... ", url)
		for _, v := range sql {
			updatesJSON, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}

			s := fmt.Sprintf(`{
		"config": %s
	}`, updatesJSON)
			post(s, url)

			time.Sleep(500 * time.Millisecond)
		}
		fmt.Println("end req... ", url)
	}
}

func post(s, reqUrl string) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(s))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("x-token", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
