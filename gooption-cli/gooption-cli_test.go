package main

import (
	"io"
	"os"
	"testing"
)

var (
	requests = []string{"Price", "Greek"}
)

type mockwriter struct{}

func (w mockwriter) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (w mockwriter) Close() error {
	return nil
}

func Test_generateTemplate(t *testing.T) {
	type args struct {
		projectConfig []config
		requests      []string
		newWriter     func(path string) (io.WriteCloser, error)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"service",
			args{
				configMap["service"],
				[]string{"Price", "Greek"},
				func(path string) (io.WriteCloser, error) { return mockwriter{}, nil },
			},
		},
		{
			"gobs",
			args{
				configMap["gobs"],
				[]string{"Price", "Greek"},
				func(path string) (io.WriteCloser, error) { return mockwriter{}, nil },
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateTemplate(tt.args.projectConfig, tt.args.requests, tt.args.newWriter)
		})
	}
}

func Test_generateJSON(t *testing.T) {
	generateJSON([]string{"ImpliedVol"}, func(path string) (io.WriteCloser, error) { return mockwriter{}, nil })
}

func Test_impliedVolRequestGenerator_generate(t *testing.T) {
	request, err := NewRequest("impliedvol")
	if err != nil {
		t.Error(err)
	}
	t.Log(request)
}

// func Test_priceRequestGenerator_generate(t *testing.T) {
// 	msg := []byte{52, 48, 52, 32, 112, 97, 103, 101, 32, 110, 111, 116, 32, 102, 111, 117, 110, 100}
// 	response := &pb.PriceResponse{}
// 	if err := proto.Unmarshal(msg, response); err != nil {
// 		t.Error(err)
// 	}
// }

// func Test_LoadJSONRequest(t *testing.T) {
// 	request := &pb.PriceRequest{}
// 	mockRequestObj, err := LoadJSONRequest("./testdata/price-handler.json", request)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(mockRequestObj)
// }
