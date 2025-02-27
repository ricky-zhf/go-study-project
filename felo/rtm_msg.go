package felo

import (
	"context"
	"github.com/sirupsen/logrus"
	"gitlab.changyinlive.com/yuyin_code_base/common_sdk/ugo"
	"gitlab.changyinlive.com/yuyin_code_base/proto_clients/room_msg/room_msg"
	"gitlab.changyinlive.com/yuyin_code_base/proto_clients/room_msg/room_msg_client"
	"time"
)

// PublicChatMsgBuffer 公屏消息缓冲区
type PublicChatMsgBuffer struct {
	msgChan   chan *room_msg.RpcSendPublicChatMsgReq // 待发送的公屏消息
	size      int                                    // 缓冲区大小
	flushTime time.Duration                          // 超时处罚
	RoomMsg   *room_msg_client.RoomMsgClient         // room_msg 客户端
}

func NewChannelAggregator(bufferSize int, flushTimeout int, client *room_msg_client.RoomMsgClient) *PublicChatMsgBuffer {
	b := &PublicChatMsgBuffer{
		msgChan:   make(chan *room_msg.RpcSendPublicChatMsgReq, bufferSize),
		size:      bufferSize,
		flushTime: time.Duration(flushTimeout) * time.Second,
		RoomMsg:   client,
	}

	// 启动后台 flush 处理
	ugo.SafeGo(func(ctx context.Context) {
		b.startFlushWorker()
	})
	return b
}

// 定期检查缓冲区并触发 flush
func (b *PublicChatMsgBuffer) startFlushWorker() {
	ticker := time.NewTicker(b.flushTime)
	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			//logrus.Debugf("ticker send public msg")
			b.FlushBuffer()

		default:
			if b.isFull() { // 缓冲区满
				logrus.Debugf("buffer full send public msg")
				b.FlushBuffer()
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// FlushBuffer 发送消息
func (b *PublicChatMsgBuffer) FlushBuffer() {
	var messages []*room_msg.RpcSendPublicChatMsgReq

	// 读取消息
	for len(b.msgChan) > 0 {
		msg := <-b.msgChan
		messages = append(messages, msg)
		if len(messages) >= b.size {
			break
		}
	}

	// 如果没有消息，则不处理
	if len(messages) == 0 {
		return
	}

	if _, err := b.RoomMsg.RpcBatchSendPublicChatMsg(context.Background(), &room_msg.RpcBatchSendPublicChatMsgReq{Reqs: messages}); err != nil {
		logrus.Errorf("RpcBatchSendPublicChatMsg failed. %s", err)
		return
	}
	logrus.Infof("RpcBatchSendPublicChatMsg success. len:%v", len(messages))
}

// PushPublicMessage 添加公屏消息到缓冲区
func (b *PublicChatMsgBuffer) PushPublicMessage(msg ...*room_msg.RpcSendPublicChatMsgReq) {
	ugo.SafeGo(func(ctx context.Context) {
		if b.isFull() {
			logrus.Debugf("Channel is full, send all msgs directly")
			if _, err := b.RoomMsg.RpcBatchSendPublicChatMsg(context.Background(), &room_msg.RpcBatchSendPublicChatMsgReq{Reqs: msg}); err != nil {
				logrus.Errorf("RpcBatchSendPublicChatMsg failed. %s", err)
			}
			return
		}

		// 放入缓冲区
		for i, v := range msg {
			select {

			case b.msgChan <- v:

			default:
				// 缓冲区已满，将剩余放进缓冲区
				newMsg := msg[i:]
				logrus.Debugf("Channel is full, send msg directly, i:%v, len:%v", i, len(newMsg))
				if _, err := b.RoomMsg.RpcBatchSendPublicChatMsg(context.Background(), &room_msg.RpcBatchSendPublicChatMsgReq{Reqs: newMsg}); err != nil {
					logrus.Errorf("RpcBatchSendPublicChatMsg failed. %s", err)
				}
				return
			}
		}
	})

}

func (b *PublicChatMsgBuffer) isFull() bool {
	return len(b.msgChan) >= b.size
}
