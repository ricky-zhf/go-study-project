package timex

import (
	"fmt"
	"time"
)

var timeLayout = "2006-01-02 15:04:05"

func Test() {
	now := time.Now()
	fmt.Println("now=", now) // 2023-09-04 10:47:10.094718 +0800 CST m=+0.000086501

	loc, _ := time.LoadLocation("UTC")
	fmt.Println("now=", now.In(loc)) // 2023-09-04 02:47:10.094718 +0000 UTC

	utc := now.UTC()
	fmt.Println("utc=", utc) // 2023-09-04 02:47:10.094718 +0000 UTC

	loc, _ = time.LoadLocation("America/New_York")
	fmt.Println("new york=", now.In(loc)) // 2023-09-03 22:47:10.094718 -0400 EDT

	// 将utc时间转换为东八区
	AsiaLocation, _ := time.LoadLocation("Asia/Shanghai")
	fmt.Println("东八区=", utc.In(AsiaLocation))

	fmt.Println(utc.In(AsiaLocation).Format(time.RFC3339))

}

// ParseUTCInLocal 将时间转换为指定时区时间
func ParseUTCInLocal(t time.Time, local string, format string) (string, error) {
	location, err := time.LoadLocation(local)
	if err != nil {
		return t.UTC().Format(format), err
	}

	return t.In(location).Format(format), nil
}
