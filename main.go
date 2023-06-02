package main

import (
	"fmt"
	_ "net/http/pprof"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var datas []string

func main() {
	now, _ := time.Parse("2006-01-02", "2023-05-31")
	//now := time.Now()
	zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	zeroTime = zeroTime.AddDate(0, 1, 1).Add(-time.Second)
	fmt.Println("当天零点时间:", zeroTime)
	//fmt.Println("当天零点时间戳:", zeroTimestamp)
}

func intSliceToString(intSlice []int) string {
	stringSlice := make([]string, len(intSlice))

	for i, v := range intSlice {
		stringSlice[i] = strconv.Itoa(v)
	}
	return strings.Join(stringSlice, ",")
}

//ParseStr 解析出string中${}和$[]内的内容，并用m对应的值替换
func ParseStr(str string, m map[string]interface{}) string {
	res := regexp.MustCompile(`\${(.*?)}`).FindAllStringSubmatch(str, -1)
	for _, match := range res {
		if v, ok := m[match[1]]; ok && v != nil {
			str = strings.ReplaceAll(str, match[0], fmt.Sprintf("%v", v))
		} else {
			str = strings.ReplaceAll(str, match[0], "")
		}
	}

	res = regexp.MustCompile(`\$\[(.*?)\]`).FindAllStringSubmatch(str, -1)
	for _, match := range res {
		if v, ok := m[match[1]]; ok && v != nil {
			str = strings.ReplaceAll(str, match[0], fmt.Sprintf("%v", v))
		} else {
			str = strings.ReplaceAll(str, match[0], "")
		}
	}
	str = regexp.MustCompile(`\s+`).ReplaceAllString(str, " ")
	return strings.TrimSpace(str)
}
