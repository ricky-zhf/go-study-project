package felo

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const (
	ProdUrl = "https://api.felolive.com"
)

func genSign(secretKey, body string) string {
	// 拼接请求体和密钥
	baseStr := body + secretKey

	// 检查基础字符串长度是否小于 20
	if len(baseStr) < 20 {
		return ""
	}

	// 分割并重排字符串
	prefix := baseStr[:10] // 前10字符
	fmt.Println("prefix:", prefix)
	postfix := baseStr[len(baseStr)-10:] // 后10字符
	fmt.Println("postfix:", postfix)
	middle := baseStr[10 : len(baseStr)-10] // 中间部分
	fmt.Println("middle:", middle)
	baseStr = postfix + middle + prefix // 重排后的baseStr
	fmt.Println("====")
	fmt.Println(baseStr)
	fmt.Println("====")

	// 生成 MD5 哈希
	hash := md5.Sum([]byte(baseStr))
	return hex.EncodeToString(hash[:])
}
