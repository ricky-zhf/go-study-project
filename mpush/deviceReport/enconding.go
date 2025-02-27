package deviceReport

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

//go:generate protoc -I. --go_out=paths=source_relative:. deviceReport.proto

func (x *DeviceReportMessage) Marshal() ([]byte, error) {
	if x == nil {
		return nil, errors.New("x is empty")
	}
	return proto.Marshal(x)
}

func (x *DeviceReportMessage) Unmarshal(data []byte) error {
	return proto.Unmarshal(data, x)
}
