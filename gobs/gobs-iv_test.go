package gooption

import (
	"reflect"
	"testing"

	"github.com/gooption/pb"
)

func Test_newOptionQuoteSliceIterator(t *testing.T) {
	type args struct {
		quotes *pb.OptionQuoteSlice
		market *pb.OptionMarket
	}
	tests := []struct {
		name string
		args args
		want *optionQuoteSliceIterator
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newOptionQuoteSliceIterator(tt.args.quotes, tt.args.market); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newOptionQuoteSliceIterator() = %v, want %v", got, tt.want)
			}
		})
	}
}
