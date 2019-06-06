package model

import (
	json "github.com/json-iterator/go"

	"github.com/golang/protobuf/proto"
)

func PB2JSON(pb proto.Message) ([]byte, error) {
	return json.Marshal(pb)
}

func JSON2Map(jsonByte []byte) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	err := json.Unmarshal(jsonByte, &data)
	return data, err
}

func InterfaceToPB(from interface{}, pb interface{}) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, pb)
}

func PB2Map(pb proto.Message) (map[string]interface{}, error) {
	if b, err := PB2JSON(pb); err != nil {
		return nil, err
	} else {
		return JSON2Map(b)
	}
}
