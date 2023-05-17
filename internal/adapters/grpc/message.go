package grpc

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/dynamic"
)

type Message struct {
	*dynamic.Message
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return m.MarshalJSONPB(&jsonpb.Marshaler{OrigName: true})
}
