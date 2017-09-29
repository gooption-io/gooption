

// /* ----------------------------------------------------------------------------
//  * This file was automatically generated by gooption-cli.
//  * 
//  * This file contains test for each service and its associated request passed to gooption-cli
//  * Do not make changes to this file unless you know what you are doing--modify the gooption-cli
//  * tmpl file instead.
//  * ----------------------------------------------------------------------------- */

package main

import (
	context "golang.org/x/net/context"

	"os"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gooption/gobs/pb"
)

var (
	s = &server{}
)


func Test_Price(t *testing.T) {
	if file, err := os.Open("./testdata/PriceRequest.json"); err == nil {
		defer file.Close()
		request := &pb.PriceRequest{}
		if jsonpb.Unmarshal(file, request) == nil {
			if response, err := s.Price(context.Background(), request); err != nil {
				t.Log(response)
				t.Error(err)
				if response.Error != "" {
					t.Error(response.Error)
				}
			}
		} else {
			t.Error(err)			
		}
	} else {
		t.Error(err)					
	}
}

func Test_Greek(t *testing.T) {
	if file, err := os.Open("./testdata/GreekRequest.json"); err == nil {
		defer file.Close()
		request := &pb.GreekRequest{}
		if jsonpb.Unmarshal(file, request) == nil {
			if response, err := s.Greek(context.Background(), request); err != nil {
				t.Log(response)
				t.Error(err)
				if response.Error != "" {
					t.Error(response.Error)
				}
			}
		} else {
			t.Error(err)			
		}
	} else {
		t.Error(err)					
	}
}

func Test_ImpliedVol(t *testing.T) {
	if file, err := os.Open("./testdata/ImpliedVolRequest.json"); err == nil {
		defer file.Close()
		request := &pb.ImpliedVolRequest{}
		if jsonpb.Unmarshal(file, request) == nil {
			if response, err := s.ImpliedVol(context.Background(), request); err != nil {
				t.Log(response)
				t.Error(err)
				if response.Error != "" {
					t.Error(response.Error)
				}
			}
		} else {
			t.Error(err)			
		}
	} else {
		t.Error(err)					
	}
}

