package asynq

import (
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

const redisAddr = "43.139.229.95:6379"
const redisAddr2 = "127.0.0.1:6379"

func run() {
	go run1()
	run2()
}

func run1() {
	srv := asynq.NewServer(
		//asynq.RedisClientOpt{Addr: redisAddr, Password: "Z3Y5V57oDcbyhouuSiRdaKph"},
		asynq.RedisClientOpt{Addr: redisAddr2},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"coupon-queue": 10,
			},

			// See the godoc for other configuration options
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				//s := int(math.Pow(float64(n), 2)) + 5
				fmt.Println("重试：", n, "间隔：")
				////fmt.Println(time.Duration(s) * time.Second)
				////return asynq.DefaultRetryDelayFunc(n, e, t)
				return 1 * time.Second
			},
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)
	mux.Handle(TypeImageResize, NewImageProcessor())
	// ...register other handlers...
	fmt.Println("run1...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func run2() {
	srv := asynq.NewServer(
		//asynq.RedisClientOpt{Addr: redisAddr, Password: "Z3Y5V57oDcbyhouuSiRdaKph"},
		asynq.RedisClientOpt{Addr: redisAddr2},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"fulfill-queue": 10,
			},

			// See the godoc for other configuration options
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				//s := int(math.Pow(float64(n), 2)) + 5
				fmt.Println("重试：", n, "间隔：")
				////fmt.Println(time.Duration(s) * time.Second)
				////return asynq.DefaultRetryDelayFunc(n, e, t)
				return 1 * time.Second
			},
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTaskV2)
	fmt.Println("run2...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
