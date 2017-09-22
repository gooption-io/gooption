package main

import (
	"reflect"
	"testing"
)

func Test_readCSV(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		{"Sample Options", args{""}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readCSV(tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("readCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}
