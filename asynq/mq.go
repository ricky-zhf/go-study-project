package asynq

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
)

// 任务处理
func handler(ctx context.Context, t *asynq.Task) error {
	switch t.Type() {
	//case TASK_EMAIL_WELCOME:
	//	var p EmailTaskPayload
	//	if err := json.Unmarshal(t.Payload(), &p); err != nil {
	//		return err
	//	}
	//case TASK_EMAIL_REMINDER:
	//	var p EmailTaskPayload
	//	if err := json.Unmarshal(t.Payload(), &p); err != nil {
	//		return err
	//	}
	default:
		return fmt.Errorf("unexpected task type: %s", t.Type())
	}
	return nil
}

//func run() {
//	srv := asynq.NewServer(
//		asynq.RedisClientOpt{Addr: "localhost:6379"},
//		asynq.Config{Concurrency: 10},
//	)
//
//	// 需要使用asynq.HandlerFunc适配
//	if err := srv.Run(asynq.HandlerFunc(handler)); err != nil {
//		log.Fatal(err)
//	}
//}
