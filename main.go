package main

import (
	"fmt"
	_ "net/http/pprof"
	"regexp"
	"strings"
)

var datas []string

func main() {
	var m map[string]string
	m["fwef"] = "fwefwe"

	fmt.Println(m["ffwef"])

	str := `{state:{NoUnkeyedLiterals:{} DoNotCompare:[] DoNotCopy:[] atomicMessageInfo:0xc0000302c0} sizeCache:0 unknownFields:[] MsgId:33010000010064821682232416061 Uuid:3301000001006482 ProductCode:NVR EventTime:1682232416 Data:fields:{key:\"alarm_event\"  value:{struct_value:{fields:{key:\"channel\"  value:{number_value:0}}  fields:{key:\"event_id\"  value:{string_value:\"1682232416\"}}  fields:{key:\"event_start\"  value:{string_value:\"1682232416\"}}  fields:{key:\"event_type\"  value:{number_value:8}}  fields:{key:\"report_type\"  value:{number_value:1}}}}} Type:DEVICE_REPORT AppId:121}`
	fmt.Println(len([]byte(str)))
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
