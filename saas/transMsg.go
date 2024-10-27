package saas

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Message struct {
	ID      int
	Content string
}

func AddMessage(i int, ctx context.Context) chan<- error {
	// 模拟添加消息操作

	ch := make(chan error, 1)
	go func(i int) {
		select {
		case err := <-ch:
			log.Println("====err,", err, ",,i=", i)
		case <-ctx.Done():
			log.Println("ctx done, i=", i)
		}
	}(i)
	return ch
}

func UpdateDatabase(i int) error {
	// 模拟数据库更新操作
	time.Sleep(1 * time.Second)

	// 假设更新失败，发送错误消息
	if i == 2 {
		return fmt.Errorf("Database update failed")
	}
	return nil
}

func do() {
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		// 添加消息
		ch := AddMessage(i, ctx)

		// 更新数据库
		if err := UpdateDatabase(i); err != nil {
			log.Println("update db failed.")
			ch <- err
		}
		cancel()
	}
}
