package pprof

import (
	"log"
	"net/http"
	"time"
)

var datas []string

func Start() {
	go func() {
		for {
			log.Printf("len: %d", Add("go-programming-tour-book"))
			time.Sleep(time.Millisecond * 1000)
		}
	}()

	_ = http.ListenAndServe(":6060", nil)

}

func Add(str string) int {
	data := []byte(str)
	datas = append(datas, string(data))
	return len(datas)
}
