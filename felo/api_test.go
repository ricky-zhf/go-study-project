package felo

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.changyinlive.com/yuyin_code_base/common_sdk/errorcode"
	"gitlab.changyinlive.com/yuyin_code_base/common_sdk/uprotobuf"
	"gitlab.changyinlive.com/yuyin_code_base/proto_clients/assets/assets"
	"gitlab.changyinlive.com/yuyin_code_base/proto_clients/assets/assets_client"
	"gitlab.changyinlive.com/yuyin_code_base/proto_clients/assets/assets_utils"
	"gitlab.changyinlive.com/yuyin_code_base/proto_clients/room_msg/room_msg"
	"google.golang.org/grpc/status"
	"runtime"
	"testing"
)

func Test_checkBalance(t *testing.T) {
	resp, err := assets_client.NewAssetsClient("felo-higo-adapter.test01.felolive.com:8081").RpcCheckUserBalance(context.Background(), &assets.RpcCheckUserBalanceReq{
		UserId:      217027933,
		AssetType:   "coin",
		AssetAmount: 1,
	})
	if err != nil && !assets_utils.IsInsufficientBalanceError(err) {
		t.Errorf("Check user balance failed. err:%v", err)
		return
	}
	fmt.Println(resp)
	// 有余额
	if assets_utils.IsInsufficientBalanceError(err) {
		fmt.Println("abcd")
	}
	fmt.Println()
}

// FromError try to convert an error to *Error.
func FromError(err error) errorcode.ErrorCode {
	if err == nil {
		return errorcode.CodeSuccess
	}

	var (
		causeErr = errors.Cause(err)
		e        errorcode.ErrorCode
	)
	if errors.As(causeErr, &e) {
		return e
	} else if grpcStatus, ok := status.FromError(causeErr); ok {
		if grpcCode := uint32(grpcStatus.Code()); grpcCode != 0 {
			return errorcode.NewErrorCode(int(grpcCode), grpcStatus.Message())
		}
	}

	return errorcode.CodeServerError
}

func Test_bonus(t *testing.T) {
	pro := getProductConfig("abc")
	fmt.Println(pro)
}

//func Test_verifyApplePayReceiptWithoutOrderid(t *testing.T) {
//	d := "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"
//	verifyApplePayReceiptWithoutOrderid(d, "2000000745557724")
//}

func Test_run(t *testing.T) {
	//tkdo()
	//loc.GetLoc()
	//run()

	//parseSignedPayload()
	testAppleRefund()

	//googleRefund()
}

func Test_sign(t *testing.T) {
	protoStruct, err := uprotobuf.ProtoMessageToProtoStruct(&room_msg.MsgContentGiftResult{
		GiftId:    "123",
		GiftClass: "123",
		//UserId:        123,
		Balance:       23.33,
		LotteryResult: nil,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(protoStruct)
}

func Test_dingding(t *testing.T) {
	client := NewDingClient("https://oapi.dingtalk.com/robot/send?access_token=33761f14a8cbfaa0b28337b87a26b5a8478af19260328e50bb1f23e5be4d7e68")

	defer func() {
		if p := recover(); p != nil {
			buf := make([]byte, 10240)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			param := &DingDingMsgSt{
				MsgType: "markdown",
			}

			content := "<font color=#FF0000 size=3 face='黑体'>【线上panic告警，请及时处理】</font>" + "\n\n**工程名称**: " + "test_project" + "\n\n**日志内容**: " + string(buf)
			param.Markdown.Title = "panic告警测试"
			param.Markdown.Text = content
			if err := DingDingSendMsg(param, client.url); err != nil {
				t.Errorf("err:%v", err)
			}
		}
	}()

	var sli []int
	sli[0] = 1
}
