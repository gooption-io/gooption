package main

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/gooption/pb"
)

func Test_impliedVolRequestGenerator_generate(t *testing.T) {
	request, err := NewRequest("impliedvol")
	if err != nil {
		t.Error(err)
	}
	t.Log(request)
}

func Test_priceRequestGenerator_generate(t *testing.T) {
	msg := []byte{52, 48, 52, 32, 112, 97, 103, 101, 32, 110, 111, 116, 32, 102, 111, 117, 110, 100}
	response := &pb.PriceResponse{}
	if err := proto.Unmarshal(msg, response); err != nil {
		t.Error(err)
	}
}

// func Test_LoadJSONRequest(t *testing.T) {
// 	request := &pb.PriceRequest{}
// 	mockRequestObj, err := LoadJSONRequest("./testdata/price-handler.json", request)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(mockRequestObj)
// }
