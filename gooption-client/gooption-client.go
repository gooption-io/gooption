//go:generate gooption-cli -p client -r Price -r Greek -r ImpliedVol
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	context "golang.org/x/net/context"

	"github.com/lehajam/dgo"
	"github.com/lehajam/dgo/protos/api"

	"github.com/lehajam/gooption/gobs/pb"
	"google.golang.org/grpc"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"

	q "github.com/lehajam/gooption/query"
	"github.com/sirupsen/logrus"
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
		conn, err = grpc.Dial(":9082", grpc.WithInsecure())
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

func query(queryString string, variables map[string]string) *api.Response {
	clientDir := clientDir()
	defer os.RemoveAll(clientDir)

	conn := dial("dgraph")
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	ctx := context.Background()
	resp, err := dg.NewTxn().QueryWithVars(ctx, queryString, variables)
	if err != nil {
		logrus.Errorln(err)
		panic(err)
	}

	logrus.Info(resp)
	return resp
}

func priceRequest() {
	resp := query(
		q.PriceRequest,
		map[string]string{
			"$timestamp":    "1514162664",
			"$optionTicker": "AAPL DEC2017 PUT",
			"$rateTicker":   "USD.FEDFUND",
		})

	priceReq := &pb.PriceRequest{}
	err := dgo.Unmarshal(resp.GetJson(), priceReq)
	if err != nil {
		panic(err)
	}

	priceReq.Pricingdate = float64(1514162664)
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
	resp := query(
		q.ImpliedvolRequest,
		map[string]string{
			"$timestamp":  "1513551151",
			"$undTicker":  "AAPL",
			"$rateTicker": "USD.FEDFUND",
		})

	ivReq := &pb.ImpliedVolRequest{}
	err := dgo.Unmarshal(resp.GetJson(), ivReq)
	if err != nil {
		panic(err)
	}
	ivReq.Pricingdate = float64(1513551151)
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

func insert(request interface{}, schema string) *api.Assigned {
	clientDir := clientDir()
	defer os.RemoveAll(clientDir)

	conn := dial("dgraph")
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)
	ctx := context.Background()

	op := &api.Operation{
		Schema: schema,
	}
	err := dg.Alter(ctx, op)
	if err != nil {
		logrus.Errorln(err)
		panic(err)
	}

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(request)
	if err != nil {
		logrus.Errorln(err)
		panic(err)
	}

	mu.SetJson = pb
	assigned, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		logrus.Errorln(err)
		panic(err)
	}

	return assigned
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

	insert(
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

	insert(
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

	insert(
		request,
		`timestamp: float @index(float) .
		ticker: string @index(exact, term) .`,
	)
}
