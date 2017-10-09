//go:generate gooption-gen client ImpliedVol
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	context "golang.org/x/net/context"

	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos"
	"github.com/gooption/gobs/pb"
	"google.golang.org/grpc"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

var (
	samplePath = "./Sample_2015_October/"
)

func dial(service string) (*grpc.ClientConn, error) {
	if service == "gobs" {
		return grpc.Dial(":50051", grpc.WithInsecure())
	} else if service == "dgraph" {
		return grpc.Dial(":9080", grpc.WithInsecure())
	} else {
		return nil, errors.New("Unknown service")
	}
}

func main() {
	clientDir, err := ioutil.TempDir("", "client_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(clientDir)

	conn1, err := dial("dgraph")
	if err != nil {
		log.Fatal(err)
	}
	defer conn1.Close()

	dgraphClient := client.NewDgraphClient([]*grpc.ClientConn{conn1}, client.DefaultOptions, clientDir)
	defer dgraphClient.Close()

	req := client.Req{}
	req.SetQuery(`{
			contract(func: eq(ticker, "AAPL DEC2017 PUT")){
			  ticker
			  strike
			  expiry
			  putcall
			}
			marketdata(func: eq(timestamp, 1507072500)) {
				spot {
					index {
						value
					}
				}
				vol {
					index {
						value
					}
				}
				rate {
					index {
						value
					}
				}
			}
		  }`)

	resp, err := dgraphClient.Run(context.Background(), &req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	printNode(0, resp.N[0])
	fmt.Printf("%+v\n", proto.MarshalTextString(resp))

	priceReq := &pb.PriceRequest{}
	err = client.Unmarshal(resp.N, priceReq)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	priceReq.Pricingdate = priceReq.Marketdata.Timestamp
	priceReq.Contract.Expiry = 1509754500
	priceReq.Contract.Putcall = pb.OptionType_CALL
	fmt.Printf("%+v\n", proto.MarshalTextString(priceReq))

	conn2, err := dial("gobs")
	if err != nil {
		log.Fatal(err)
	}
	defer conn2.Close()

	gobsClient := pb.NewGoBSServerClient(conn2)
	price, err := gobsClient.Price(context.Background(), priceReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("price: %f\n", price.Price)
}

func printNode(depth int, node *protos.Node) {

	fmt.Println(strings.Repeat(" ", depth), "Atrribute : ", node.Attribute)

	// the values at this level
	for _, prop := range node.GetProperties() {
		fmt.Println(strings.Repeat(" ", depth), "Prop : ", prop.Prop, " Value : ", prop.Value, " Type : %T", prop.Value)
	}

	for _, child := range node.Children {
		fmt.Println(strings.Repeat(" ", depth), "+")
		printNode(depth+1, child)
	}
}

func insertPriceRequest() {
	if file, err := os.Open("./testdata/PriceRequest.json"); err == nil {
		defer file.Close()
		request := &pb.PriceRequest{}
		if jsonpb.Unmarshal(file, request) == nil {
			insertRequest(request)
		} else {
			fmt.Printf("%s", err.Error())
		}
	} else {
		fmt.Printf("%s", err.Error())
	}
}

func insertGreekRequest() {
	if file, err := os.Open("./testdata/GreekRequest.json"); err == nil {
		defer file.Close()
		request := &pb.GreekRequest{}
		if jsonpb.Unmarshal(file, request) == nil {
			insertRequest(request)
		} else {
			fmt.Printf("%s", err.Error())
		}
	} else {
		fmt.Printf("%s", err.Error())
	}
}

func insertImpliedVolRequest() {
	if file, err := os.Open("./testdata/ImpliedVolRequest.json"); err == nil {
		defer file.Close()
		request := &pb.ImpliedVolRequest{}
		if jsonpb.Unmarshal(file, request) == nil {
			insertRequest(request)
		} else {
			fmt.Printf("%s", err.Error())
		}
	} else {
		fmt.Printf("%s", err.Error())
	}
}

func insertRequest(request interface{}) {
	clientDir, err := ioutil.TempDir("", "client_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(clientDir)

	conn1, err := dial("dgraph")
	if err != nil {
		log.Fatal(err)
	}
	defer conn1.Close()

	dgraphClient := client.NewDgraphClient([]*grpc.ClientConn{conn1}, client.DefaultOptions, clientDir)
	defer dgraphClient.Close()

	req := client.Req{}
	req.SetObject(request)
	resp, err := dgraphClient.Run(context.Background(), &req)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%+v\n", proto.MarshalTextString(resp))

	req.SetObject(request)
	resp, err = dgraphClient.Run(context.Background(), &req)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%+v\n", proto.MarshalTextString(resp))
}
