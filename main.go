package main

import (
	"StudyProject/mpush"
	"fmt"
	"net/http"
)

func methodA(w http.ResponseWriter, r *http.Request) {
	for i := 1; i <= 100; i++ {
		fmt.Println(i)
	}
}

func main() {
	mpush.Do()
}
