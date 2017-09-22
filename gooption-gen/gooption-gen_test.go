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
		projectConfig []ArgT
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
				projectMap["service"],
				[]string{"Price", "Greek"},
				func(path string) (io.WriteCloser, error) { return mockwriter{}, nil },
			},
		},
		{
			"gobs",
			args{
				projectMap["gobs"],
				[]string{"Price", "Greek"},
				func(path string) (io.WriteCloser, error) { return mockwriter{}, nil },
			},
		},
		{
			"utils",
			args{
				[]ArgT{{TemplateName: "request_utils", FileName: "request_utils.go"}},
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
