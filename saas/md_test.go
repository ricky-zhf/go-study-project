package saas

import (
	"fmt"
	"testing"
)

type AppHeader struct {
	Appid      string `header:"appid" binding:"required"`
	Lang       string `header:"lang" binding:"required"`
	Timestamp  string `header:"timestamp" binding:"required"`
	Nonce      string `header:"nonce" binding:"required"`
	Source     string `header:"source" binding:"required"`
	Sign       string `header:"sign" binding:"required"`
	Token      string `header:"token"`
	AppVersion string `header:"app_version"`
}

func TestMd5Encode(t *testing.T) {
	header := AppHeader{}
	header.AppVersion = "215"
	header.Appid = "60A22968B9A63"
	header.Lang = "cn"
	header.Nonce = "4483837303306280"
	header.Source = "harmoryOS"
	header.Timestamp = "1715829490"
	//header.Token = "b22333edaac464f0bc8ecc82f3561d1f"

	sign := "dc3cf642cca6d15554b19432ef186d3e"

	str := ""
	if len(header.AppVersion) > 0 {
		str = header.AppVersion + header.Appid + header.Lang + header.Nonce + header.Source + header.Timestamp
	} else {
		str = header.Appid + header.Lang + header.Nonce + header.Source + header.Timestamp
	}
	if header.Token != "" {
		str += header.Token
	}

	encode := Md5Encode(str)
	fmt.Println("解析:", encode)
	fmt.Println("传参:", sign)
}

func TestAdminUpdate(t *testing.T) {
	//run()

	do()
}
