package asynq

import (
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"testing"
)

func Test_server_run(t *testing.T) {
	run()
}

func Test_client_run(t *testing.T) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr, Password: "Z3Y5V57oDcbyhouuSiRdaKph"})
	//client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr2})

	defer client.Close()
	userCoupon(client)

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	//email(client)

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

func email(client *asynq.Client) {
	task, err := NewEmailDeliveryTask(42, "9999")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task, asynq.MaxRetry(3), asynq.Queue("fulfill-queue"))
	info, err = client.Enqueue(task, asynq.MaxRetry(3), asynq.Queue("coupon-queue"))

	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}

func userCoupon(client *asynq.Client) {
	payload, err := json.Marshal(UseCouponRequest{
		UserId:       99,
		Uuid:         "43534",
		PackageId:    345,
		CouponId:     534,
		OrderNo:      "3453434",
		PackagePrice: 1123,
	})
	if err != nil {
		fmt.Println(err)
	}
	task := asynq.NewTask("LABJ-Task-UseCoupon", payload)
	info, err := client.Enqueue(task, asynq.Queue("LABJ-Queue-UseCoupon"), asynq.MaxRetry(5))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}


172.20.12.2:2379

134.175.211.197:2379

var mobilePushConfig model.MobilePushConfig
config, err := l.svcCtx.NotificationConfigModel.FindOneByAppIdNotificationTypeIsDefault(l.ctx, user.AppId, model.NotificationTypeMobilePush, model.DefaultNotificationConfig)
if err != nil {
l.Errorf("Find notification config failed. error=%v", err)
return err
}
if err = json.Unmarshal([]byte(config.Configuration), &mobilePushConfig); err != nil {
l.Errorf("Unmarshal notification failed. err=%v", err)
return err
}

pushPem, err := l.svcCtx.NotificationPushPemModel.FindOneByFilename(l.ctx, mobilePushConfig.AppleKeyFile)
if err != nil {
l.Errorf("Find notification push pem failed. error=%v", err)
return err
}
mobilePushConfig.AppleValFileValue = pushPem.Content
