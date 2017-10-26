//go:generate gooption-gen client ImpliedVol
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	context "golang.org/x/net/context"

	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos"
	"github.com/gooption/gobs/pb"
	"google.golang.org/grpc"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

func main() {
}

func dial(service string) *grpc.ClientConn {
	var (
		conn *grpc.ClientConn
		err  error
	)

	if service == "gobs" {
		conn, err = grpc.Dial(":50051", grpc.WithInsecure())
	} else if service == "dgraph" {
		conn, err = grpc.Dial(":9080", grpc.WithInsecure())
	} else {
		err = errors.New("Unknown service")
	}

	if err != nil {
		log.Fatal(service)
		log.Fatal(err)
		panic(service + "\n" + err.Error())
	}

	return conn
}

func clientDir() string {
	dir, err := ioutil.TempDir("", "client_")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return dir
}

func priceRequest() {
	clientDir := clientDir()
	defer os.RemoveAll(clientDir)

	conn1 := dial("dgraph")
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
	
			marketdata(func: eq(timestamp, 1508274400)) @cascade { 
				spot {
					...indexInfo
				}
				vol  {
					...indexInfo
				}
				rate  {
					...indexInfo
				}
			} 
		}
			
		fragment indexInfo {
			index @filter(eq(ticker, "AAPL") or eq(ticker, "USD.FEDFUND")) {
				timestamp
				ticker
				value
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
	priceReq.Pricingdate = float64(time.Now().Unix())
	priceReq.Contract.Putcall = pb.OptionType_CALL
	fmt.Printf("%+v\n", proto.MarshalTextString(priceReq))

	conn2 := dial("gobs")
	defer conn2.Close()

	gobsClient := pb.NewGobsClient(conn2)
	price, err := gobsClient.Price(context.Background(), priceReq)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("price: %f\n", price.Price)
}

func ivRequest() {
	clientDir := clientDir()
	defer os.RemoveAll(clientDir)

	conn1 := dial("dgraph")
	defer conn1.Close()

	dgraphClient := client.NewDgraphClient([]*grpc.ClientConn{conn1}, client.DefaultOptions, clientDir)
	defer dgraphClient.Close()

	req := client.Req{}
	req.SetQuery(`
		{
			marketdata(func: eq(timestamp, 1508274400)) @cascade { 
				spot {
					...indexInfo
				}
				vol  {
					...indexInfo
				}
				rate  {
					...indexInfo
				}
			} 
	
			quotes(func: eq(timestamp, 1508274400)) @cascade { 
			  expiry
			  puts {
				...quote
			  }
			  calls {
				...quote
			  }
			} 
		}
		  
		fragment quote {
			strike
			bid
			ask
			openinterest
		}
	
		fragment indexInfo {
			index @filter(eq(ticker, "AAPL") or eq(ticker, "USD.FEDFUND")) {
				timestamp
				ticker
				value
			}
		}`)

	resp, err := dgraphClient.Run(context.Background(), &req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	printNode(0, resp.N[0])
	fmt.Printf("%+v\n", proto.MarshalTextString(resp))

	ivReq := &pb.ImpliedVolRequest{}
	err = client.Unmarshal(resp.N, ivReq)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ivReq.Pricingdate = float64(time.Now().Unix())
	fmt.Printf("%+v\n", proto.MarshalTextString(ivReq))

	conn2 := dial("gobs")
	defer conn2.Close()

	gobsClient := pb.NewGobsClient(conn2)
	volSurf, err := gobsClient.ImpliedVol(context.Background(), ivReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("***************gobsClient response***************")
	fmt.Printf("%+v\n", proto.MarshalTextString(volSurf))
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

	conn1 := dial("dgraph")
	defer conn1.Close()

	dgraphClient := client.NewDgraphClient([]*grpc.ClientConn{conn1}, client.DefaultOptions, clientDir)
	defer dgraphClient.Close()

	req := client.Req{}
	req.SetObject(request)
	req.SetSchema(`
		timestamp: float @index(float) .
		ticker: string @index(exact, term) .`)
	resp, err := dgraphClient.Run(context.Background(), &req)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%+v\n", proto.MarshalTextString(resp))
}
