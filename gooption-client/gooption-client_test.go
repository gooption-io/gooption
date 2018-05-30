package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/lehajam/gooption/gobs/pb"
	"github.com/mitchellh/mapstructure"

	"github.com/sirupsen/logrus"
)

func Test_insertImpliedVolRequest(t *testing.T) {
	insertImpliedVolRequest()
}

func Test_insertPriceRequest(t *testing.T) {
	insertPriceRequest()
}

func Test_priceRequest(t *testing.T) {
	priceRequest()
}

func Test_ivRequest(t *testing.T) {
	ivRequest()
}

func Test_contractRequest(t *testing.T) {
	resp := query(
		`query ContractRequest($optionTicker: string){
			contract(func: eq(ticker, $optionTicker)){
				ticker
				strike
				undticker
				expiry
				putcall
			}
		}`,
		map[string]string{
			"$optionTicker": "AAPL DEC2017 PUT",
		})

	priceReq := &pb.PriceRequest{}

	// Fails with error "json: cannot unmarshal array into Go struct field PriceRequest.contract of type pb.European"
	err := json.Unmarshal(resp.GetJson(), priceReq)
	if err != nil {
		type Root struct {
			Contract []pb.European `json:"contract"`
		}

		root := &Root{}
		// Works fine
		err = json.Unmarshal(resp.GetJson(), root)
		if err != nil {
			t.Error(err)
		}
		t.Log(root)
		t.Fail()
	}
	t.Log(priceReq)
}

func Test_contractRequestJsonpb(t *testing.T) {
	resp := query(
		`query ContractRequest($optionTicker: string){
			contract(func: eq(ticker, $optionTicker)){
				ticker
				strike
				undticker
				expiry
				putcall
			}
		}`,
		map[string]string{
			"$optionTicker": "AAPL DEC2017 PUT",
		})

	priceReq := &pb.PriceRequest{}
	// resolver := funcResolver(func(turl string) (proto.Message, error) {
	// 	logrus.Info(turl)
	// 	return nil, nil
	// })
	unmarshaler := jsonpb.Unmarshaler{}
	err := unmarshaler.Unmarshal(bytes.NewReader(resp.GetJson()), priceReq)
	if err != nil {
		logrus.Error(err)
		type Root struct {
			Contract []pb.European `json:"contract"`
		}

		root := &Root{}
		// Works fine
		err = json.Unmarshal(resp.GetJson(), root)
		if err != nil {
			t.Error(err)
		}
		t.Log(root)
	}
	t.Log(priceReq)
}

// extractFile extracts a FileDescriptorProto from a gzip'd buffer.
func extractFile(gz []byte) (*descriptor.FileDescriptorProto, error) {
	r, err := gzip.NewReader(bytes.NewReader(gz))
	if err != nil {
		return nil, fmt.Errorf("failed to open gzip reader: %v", err)
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to uncompress descriptor: %v", err)
	}

	fd := new(descriptor.FileDescriptorProto)
	if err := proto.Unmarshal(b, fd); err != nil {
		return nil, fmt.Errorf("malformed FileDescriptorProto: %v", err)
	}

	return fd, nil
}

func Test_describeProto(t *testing.T) {
	// msg := &descriptor.FileDescriptorProto{}
	fd, err := extractFile(proto.FileDescriptor("/home/lehajam/go/src/github.com/lehajam/gooption/pb/contract.proto"))
	if err != nil {
		logrus.Error(err)
		panic(err)
	}

	// err = proto.Unmarshal(fd, msg)
	// if err != nil {
	// 	logrus.Error(err)
	// 	panic(err)
	// }
	logrus.Info(fd)

	// msg := &pb.PriceRequest{}
	// _, md := descriptor.ForMessage(msg)
	// fields := md.GetField()
	// scalars := make([]string, 0)
	// vectors := make([]string, 0)
	// for _, f := range fields {
	// 	logrus.Info(f)
	// 	if f.IsMessage() {
	// 		if f.IsRepeated() {
	// 			vectors = append(vectors, f.GetJsonName())
	// 		} else {
	// 			scalars = append(scalars, f.GetJsonName())
	// 		}
	// 	}
	// }

	// logrus.Info(scalars)
	// logrus.Info(vectors)

	// type JSONMessage struct {
	// 	Scalar bool     `json:"scalar"`
	// 	Fields []string `json:"fields"`
	// }
	// var data JSONMessage
	// if len(scalars) < len(vectors) {
	// 	data = JSONMessage{
	// 		Scalar: true,
	// 		Fields: scalars,
	// 	}
	// } else {
	// 	data = JSONMessage{
	// 		Scalar: false,
	// 		Fields: vectors,
	// 	}
	// }

	// f, err := os.Create("protofields.json")
	// if err != nil {
	// 	t.Error(err)
	// }
	// json.NewEncoder(f).Encode(data)
}

func Test_decodeJSON(t *testing.T) {
	resp := query(
		`query ContractRequest($optionTicker: string){
			contract(func: eq(ticker, $optionTicker)){
				ticker
				strike
				undticker
				expiry
				putcall
			}
		}`,
		map[string]string{
			"$optionTicker": "AAPL DEC2017 PUT",
		})

	var parsed map[string]interface{}
	err := json.Unmarshal(resp.GetJson(), &parsed)
	if err != nil {
		t.Error(err)
	}

	var result pb.PriceRequest
	config := &mapstructure.DecoderConfig{
		Result: &result,
		DecodeHook: func(
			f reflect.Kind,
			t reflect.Kind,
			data interface{}) (interface{}, error) {

			switch t {
			case reflect.Struct:
				switch f {
				case reflect.Slice:
					dataAsSlice := data.([]interface{})
					if len(dataAsSlice) == 0 {
						return nil, nil
					}
					return dataAsSlice[0], nil
				}
			}

			return data, nil
		},
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		panic(err)
	}
	if err := decoder.Decode(parsed); err != nil {
		panic(err)
	}
	logrus.Info(result)
	logrus.Info(result.Contract.Ticker)
}
