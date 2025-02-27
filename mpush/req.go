package mpush

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func main() {
	appVersion := "212"
	appId := "60A22968B9A63"
	lang := "cn"
	nonce := "ajVp1HFU696HNyeH"
	source := "Android"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	token := "59190b0f2364c0e3d065f13de03c1010"
	input := appVersion + appId + lang + nonce + source + timestamp + token
	sign := Md5Encode(input)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://app-pre.wjacloud.com/api/v1/device/detail?device_id=3301000000865917", nil)
	if err != nil {
		fmt.Println(err)
	}
	headers := http.Header{}
	headers.Add("source", source)
	headers.Add("token", token)
	headers.Add("appid", appId)
	headers.Add("app_version", appVersion)
	headers.Add("sign", sign)
	headers.Add("nonce", nonce)
	headers.Add("lang", lang)
	headers.Add("timestamp", timestamp)
	headers.Add("appversion", appVersion)

	req.Header = headers

	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf(string(body))

}

type S struct {
}

func Md5Encode(encodeString string) string {
	h := md5.New()
	h.Write([]byte(encodeString))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
