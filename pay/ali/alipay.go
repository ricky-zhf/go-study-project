package ali

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func pay() {
	url := "https://openapi.alipay.com/gateway.do?app_id=2021004119669087&biz_content=%257B%2522out_trade_no%2522:%252226819928321846476800000000000000%2522,%2522product_code%2522:%2522QUICK_WAP_WAY%2522,%2522subject%2522:%2522%25E5%25A2%259E%25E5%2580%25BC%25E6%259C%258D%25E5%258A%25A1subject%2522,%2522total_amount%2522:0.01%257D&charset=utf-8&format=JSON&method=alipay.trade.wap.pay&notify_url=https://localhost/api/v1/notify/alipay/26819928321846476800000000000000&return_url=https://insight-t.xiaozlife.com/packages/index.html"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "openapi.alipay.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", "JSESSIONID=F9A8E9F751AB83356FCE3C3138C26191; zone=GZ00C; ALIPAYJSESSIONID=GZ00w8kMJVFM6St5SlBKjILXyad4dKsuperapiGZ00; ctoken=bKd-BZzGa_88ezet; spanner=7sxvwpXj3iJc2Wc3LDA3QiG1T5Orzsr/")

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
