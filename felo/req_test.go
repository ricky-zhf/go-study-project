package felo

import (
	"gitlab.changyinlive.com/yuyin_code_base/common_sdk/uprotobuf"
	"gitlab.changyinlive.com/yuyin_code_base/proto_clients/room_msg/room_msg"
	"gitlab.changyinlive.com/yuyin_code_base/proto_clients/room_msg/room_msg_client"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_PublicChatMsgBuffer(t *testing.T) {
	aggregator := NewChannelAggregator(100, 2, room_msg_client.NewRoomMsgClient("felo-room-msg.test01.felolive.com:8081"))
	var publicMsgReq []*room_msg.RpcSendPublicChatMsgReq
	for i := 0; i < 100; i++ {
		publicContent, err := uprotobuf.ProtoMessageToProtoStruct(&room_msg.MsgContentGiftResult{
			GiftId:    "20104",
			GiftClass: "luck_gift",
		})
		if err != nil {
			return
		}
		publicMsgReq = append(publicMsgReq, &room_msg.RpcSendPublicChatMsgReq{
			FromUserId:   217027933,
			RoomId:       110023,
			ContentType:  room_msg.PublicChatMsgContentType_lucky_gift,
			Content:      publicContent,
			InteractType: room_msg.PublicChatMsgInteractType_user_msg,
			BizType:      "send_gift_room",
			CtMills:      time.Now().UnixMilli(),
		})
	}

	aggregator.PushPublicMessage(publicMsgReq...)

	select {}
}

func Test_sendReq(t *testing.T) {
	method := "POST"
	//me_account:auth:felo:

	api := "/felo/api/app/translate/v1/app_translate/get"
	userId := "100102034"
	secretKey := "ZDY0MTBlODcx"
	token := "T2KdNFBU3i149ivgM_7GqjhFTSFng0xX0jzxz4rGsQJQb13bXrXO_2dczXT4c87R85VbL"
	reqBody := `{
    "translation_keys": [
        "v110_intimacy_title",
        "v110_intimacy_Instructions_2",
        "v110_intimacy_Instructions_5",
        "v110_intimacy_Instructions_6",
        "v110_intimacy_Instructions_7"
    ],
    "base_params": {
        "x_app_type_k": "android",
        "x_package_name_k": "com.example.strawberry",
        "x_device_id_k": "6d4f4de6b1feece1",
        "x_product_name_k": "felo",
        "x_install_referer_k": "",
        "x_user_id_k": 100102034,
        "x_version_name_k": "1.0.0.560",
        "x_country_code_k": "",
        "x_device_language_k": "en",
        "x_channel_k": "test",
        "x_version_code_k": 10000
    }
}`

	payload := strings.NewReader(reqBody)
	url := ProdUrl + api + "?sign=" + genSign(secretKey, reqBody)

	t.Logf("Req url:%v, body:%v", url, reqBody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Add("X-Felo-User", userId)
	req.Header.Add("X-Felo-Token", token)

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("res:%v", string(body))
}
