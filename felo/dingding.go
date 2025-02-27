package felo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type DingClient struct {
	url string
}

func NewDingClient(url string) *DingClient {
	return &DingClient{
		url: url,
	}
}

func DingDingSendMsg(param *DingDingMsgSt, URL string) (err error) {
	buf, err := json.Marshal(param)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", URL, bytes.NewReader(buf))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{
		Timeout: time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var respData struct {
		ErrCode int64  `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return
	}
	if respData.ErrCode != 0 {
		err = fmt.Errorf("dingding resp: %s", body)
	}
	return
}

type DingDingMsgSt struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Link struct {
		Text       string `json:"text"`
		Title      string `json:"title"`
		PicURL     string `json:"picUrl"`
		MessageURL string `json:"messageUrl"`
	} `json:"link"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At struct {
		AtMobiles []string `json:"atMobiles"`
		IsAtAll   bool     `json:"isAtAll"`
	} `json:"at"`
}
