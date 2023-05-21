package main

import (
	"fmt"
	_ "net/http/pprof"
	"regexp"
	"strconv"
	"strings"
)

var datas []string

func main() {
	fmt.Println(intSliceToString([]int{1, 2, 3, 3, 3, 3, 3, 1, 2, 3, 3, 4, 4}))
}

func intSliceToString(intSlice []int) string {
	stringSlice := make([]string, len(intSlice))

	for i, v := range intSlice {
		stringSlice[i] = strconv.Itoa(v)
	}
	return strings.Join(stringSlice, ",")
}
func parseStr(str string, m map[string]interface{}) string {
	// \$\[([^\]]+)\]|\$\{([^}]+)\}
	res := regexp.MustCompile(`\$\[([^\]]+)\]|\$\{([^}]+)\}`).FindAllStringSubmatch(str, -1)
	for _, match := range res {
		if v, ok := m[match[1]]; ok {
			str = strings.ReplaceAll(str, match[0], fmt.Sprintf("%v", v))
		} else {
			str = strings.ReplaceAll(str, match[0], "")
		}
	}
	res = regexp.MustCompile(`\$\[(.*?)\]`).FindAllStringSubmatch(str, -1)
	for _, match := range res {
		if v, ok := m[match[1]]; ok {
			str = strings.ReplaceAll(str, match[0], fmt.Sprintf("%v", v))
		} else {
			str = strings.ReplaceAll(str, match[0], "")
		}
	}

	str = regexp.MustCompile(`\s+`).ReplaceAllString(str, " ")
	return strings.TrimSpace(str)
}
