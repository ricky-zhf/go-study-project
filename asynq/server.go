package asynq

import (
	"github.com/hibiken/asynq"
	"log"
	"time"
)

const redisAddr = "127.0.0.1:6379"

func run() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},

			// See the godoc for other configuration options
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				//s := int(math.Pow(float64(n), 2)) + 5
				//fmt.Println("重试：", n, "间隔：", s)
				////fmt.Println(time.Duration(s) * time.Second)
				////return asynq.DefaultRetryDelayFunc(n, e, t)
				return time.Duration(n) * time.Second
			},
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
	mux.Handle(TypeImageResize, NewImageProcessor())
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
