package asynq

import (
	"github.com/hibiken/asynq"
	"log"
	"testing"
)

func Test_server_run(t *testing.T) {
	run()
}

func Test_client_run(t *testing.T) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	task, err := NewEmailDeliveryTask(42, "9999")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task, asynq.MaxRetry(10))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------
	//task, err = NewEmailDeliveryTask(43, "some:template:id")
	//info, err = client.Enqueue(task, asynq.ProcessIn(10*time.Second))
	//if err != nil {
	//	log.Fatalf("could not schedule task: %v", err)
	//}
	//log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune task processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------

	//task, err = NewImageResizeTask("https://example.com/myassets/image.jpg")
	//if err != nil {
	//	log.Fatalf("could not create task: %v", err)
	//}
	//info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	//if err != nil {
	//	log.Fatalf("could not enqueue task: %v", err)
	//}
	//log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
