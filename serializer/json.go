package serializer

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

//ProtobufToJSON conversts protocol buffer to json string
func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       " ",
		OrigName:     true,
	}

	data, err := marshaler.MarshalToString(message)

	return string(data), err
}
