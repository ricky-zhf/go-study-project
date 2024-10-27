package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"sync"
)

type st struct {
	a string
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("watcher error")
	}
	defer watcher.Close()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Printf("event type:%s\n", event.Op)
				fmt.Printf("file name:%#v\n", event.Name)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println(err.Error())
			}
		}
		wg.Done()
	}()
	err = watcher.Add("/Users/zhouhuaifeng/GoWorkspace/src/StudyProject/README.md")
	if err != nil {
		fmt.Println(err.Error())
	}
	wg.Wait()
}
