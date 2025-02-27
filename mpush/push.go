package mpush

import (
	"StudyProject/mpush/deviceReport"
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/structpb"
	"net/http"
)

func methodA(w http.ResponseWriter, r *http.Request) {
	newStruct, _ := structpb.NewStruct(map[string]interface{}{
		"alarm_event": map[string]interface{}{
			"event_end":   "1680596000",
			"event_id":    "123",
			"event_start": "1680595978",
			"event_type":  0,
			"report_type": 1,
		},
	})

	msg := deviceReport.DeviceReportMessage{
		MsgId:       "aa123",
		Uuid:        "3301100005204779",
		ProductCode: "Qsee_Q47",
		EventTime:   1680596000,
		Data:        newStruct,
		Type:        "DEVICE_REPORT",
		AppId:       6369,
	}
	d, _ := msg.Marshal()

	//134.175.211.197:9092 dev
	//106.52.97.111:9092 test
	writer := &kafka.Writer{
		Addr:     kafka.TCP([]string{"106.52.97.111:9092"}...),
		Balancer: &kafka.Hash{},
	}

	_ = writer.WriteMessages(context.Background(), kafka.Message{
		Topic:   "topic_device_report",
		Key:     []byte(msg.Uuid),
		Value:   d,
		Headers: []kafka.Header{{Key: "trace", Value: []byte("123123123")}},
	})
	fmt.Println("push end. msg=%+v", msg.String())
}

func Do() {
	http.HandleFunc("/", methodA)
	http.ListenAndServe(":9988", nil)
}
