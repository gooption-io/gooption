//go:generate gooption-cli -p client -r Price -r Greek -r ImpliedVol
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
	"github.com/lehajam/gooption/gobs/pb"
	"google.golang.org/grpc"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"

	"github.com/lehajam/gooption/query"
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

func dgraphClient(query string) *protos.Response {
	clientDir := clientDir()
	defer os.RemoveAll(clientDir)

	conn1 := dial("dgraph")
	defer conn1.Close()

	dgraphClient := client.NewDgraphClient([]*grpc.ClientConn{conn1}, client.DefaultOptions, clientDir)
	defer dgraphClient.Close()

	req := client.Req{}
	req.SetQuery(query)

	resp, err := dgraphClient.Run(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	printNode(0, resp.N[0])
	fmt.Printf("%+v\n", proto.MarshalTextString(resp))

	return resp
}

func priceRequest() {
	resp := dgraphClient(
		query.GetPriceRequestQuery("1513551151", "AAPL DEC2017 PUT", "USD.FEDFUND"),
	)

	priceReq := &pb.PriceRequest{}
	err := client.Unmarshal(resp.N, priceReq)
	if err != nil {
		panic(err)
	}

	priceReq.Pricingdate = float64(time.Now().Unix())
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
	resp := dgraphClient(
		query.GetImpliedVolRequestQuery("1513551151", "AAPL", "USD.FEDFUND"),
	)

	ivReq := &pb.ImpliedVolRequest{}
	err := client.Unmarshal(resp.N, ivReq)
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
	file, err := os.Open("./testdata/PriceRequest.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	request := &pb.PriceRequest{}
	err = jsonpb.Unmarshal(file, request)
	if err != nil {
		panic(err)
	}

	insertRequest(
		request,
		`timestamp: float @index(float) .
		ticker: string @index(exact, term) .
		undticker: string @index(exact, term) .`,
	)
}

func insertGreekRequest() {
	file, err := os.Open("./testdata/GreekRequest.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	request := &pb.GreekRequest{}
	err = jsonpb.Unmarshal(file, request)
	if err != nil {
		panic(err)
	}

	insertRequest(
		request,
		`timestamp: float @index(float) .
		ticker: string @index(exact, term) .
		undticker: string @index(exact, term) .`,
	)
}

func insertImpliedVolRequest() {
	file, err := os.Open("./testdata/ImpliedVolRequest.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	request := &pb.ImpliedVolRequest{}
	err = jsonpb.Unmarshal(file, request)
	if err != nil {
		panic(err)
	}

	insertRequest(
		request,
		`timestamp: float @index(float) .
		ticker: string @index(exact, term) .`,
	)
}

func insertRequest(request interface{}, schema string) {
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
	req.SetSchema(schema)
	req.SetObject(request)
	resp, err := dgraphClient.Run(context.Background(), &req)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%+v\n", proto.MarshalTextString(resp))
}
