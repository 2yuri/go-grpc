package serializer

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

//WriteProtobufToJSONFile WRITE BUFF TO JSON
func WriteProtobufToJSONFile(message proto.Message, filename string) error {
	data, err := ProtobufToJSON(message)

	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary: %w", err)
	}

	err = ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("cannot save json: %w", err)
	}

	return nil
}

//WriteProtobufToBinaryFile WRITE BUFF TO BINARY
func WriteProtobufToBinaryFile(message proto.Message, filename string) error {
	data, err := proto.Marshal(message)

	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary: %w", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write file to binary: %v", err)
	}

	return nil
}

//ReadProtobufFromBinaryFile reads protocol message from binary file
func ReadProtobufFromBinaryFile(filename string, message proto.Message) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read data from file: %w", err)
	}

	error := proto.Unmarshal(data, message)
	if error != nil {
		return fmt.Errorf("cannot unmarshal binary: %w", err)
	}

	return nil
}
