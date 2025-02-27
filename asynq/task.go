package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/runtime/protoimpl"
	"log"
	"time"
)

// A list of task types.
const (
	TypeEmailDelivery = "email:deliver"
	TypeImageResize   = "image:resize"
)

type EmailDeliveryPayload struct {
	UserID     int
	TemplateID string
}

type ImageResizePayload struct {
	SourceURL string
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

// NewEmailDeliveryTask 生成发送邮件的task
func NewEmailDeliveryTask(userID int, tmplID string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{UserID: userID, TemplateID: tmplID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, payload), nil
}

// NewImageResizeTask 生成resize图像的task
func NewImageResizeTask(src string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageResizePayload{SourceURL: src})
	if err != nil {
		return nil, err
	}
	// task options can be passed to NewTask, which can be overridden at enqueue time.
	return asynq.NewTask(TypeImageResize, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

//---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input task.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
//---------------------------------------------------------------

// HandleEmailDeliveryTask 发送邮件
func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("V1====Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	// Email delivery code ...
	return errors.New("failed")
}

// HandleEmailDeliveryTaskV2 发送邮件
func HandleEmailDeliveryTaskV2(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("V2====Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	// Email delivery code ...
	return errors.New("failed")
}

// ImageProcessor implements asynq.Handler interface.
type ImageProcessor struct {
	// ... fields for structa
}

func (processor *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Resizing image: src=%s", p.SourceURL)
	// Image resizing code ...
	return nil
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

// saas test
type UseCouponRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId       int64  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`                   // 用户ID
	Uuid         string `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`                                      // 设备ID
	PackageId    int64  `protobuf:"varint,3,opt,name=package_id,json=packageId,proto3" json:"package_id,omitempty"`          // 套餐ID
	CouponId     int64  `protobuf:"varint,4,opt,name=coupon_id,json=couponId,proto3" json:"coupon_id,omitempty"`             // 优惠券ID
	OrderNo      string `protobuf:"bytes,5,opt,name=order_no,json=orderNo,proto3" json:"order_no,omitempty"`                 // 订单号
	PackagePrice int64  `protobuf:"varint,6,opt,name=package_price,json=packagePrice,proto3" json:"package_price,omitempty"` // 套餐价格 * 100
}

func NewUserCouponRequestTask() (*asynq.Task, error) {
	payload, err := json.Marshal(UseCouponRequest{
		UserId:       11111,
		Uuid:         "43534",
		PackageId:    345,
		CouponId:     534,
		OrderNo:      "3453434",
		PackagePrice: 1123,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask("LABJ-Task-UseCoupon", payload), nil
}
