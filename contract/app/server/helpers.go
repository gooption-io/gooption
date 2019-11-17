package server

import (
	api_pb "contract/api"
	"errors"
	"io/ioutil"
	"log"
	"os"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/lehajam/dgo"
	"github.com/lehajam/dgo/protos/api"

	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

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

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func newDgraphClient() *dgo.Dgraph {
	clientDir := clientDir()
	defer os.RemoveAll(clientDir)

	conn := dial("dgraph")
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	return dgo.NewDgraphClient(dc)
}

func query(ctx context.Context, db *dgo.Dgraph, queryString string, variables map[string]string) *api.Response {
	if variables == nil {
		resp, err := db.NewTxn().Query(ctx, queryString)
		if err != nil {
			logrus.Errorln(err)
			panic(err)
		}

		return resp
	} else {
		resp, err := db.NewTxn().QueryWithVars(ctx, queryString, variables)
		if err != nil {
			logrus.Errorln(err)
			panic(err)
		}

		return resp
	}
}

func getAllContracts(ctx context.Context, db *dgo.Dgraph) []*api_pb.European {
	query := `{
		time(func: has(calls), orderdesc: timestamp, first: 1) {
			T as timestamp
		}

		contracts(func: has(~calls), first: 50) @filter(eq(timestamp, val(T))) {
			uid
			strike
			putcall
			ask
			~calls { expiry }
		}
	}`

	res := &struct {
		Time []struct {
			Timestamp float64 `json:"timestamp"`
		} `json:"time"`

		Contracts []struct {
			Ask     float64 `json:"ask"`
			Putcall string  `json:"putcall"`
			Strike  float64 `json:"strike"`
			UID     string  `json:"uid"`
			Calls   []struct {
				Expiry float64 `json:"expiry"`
			} `json:"~calls"`
		} `json:"contracts"`
	}{}

	dRes, err := db.NewTxn().Query(ctx, query)
	if err != nil {
		logrus.Errorln(err)
		panic(err)
	}

	json.Unmarshal(dRes.GetJson(), &res)
	if err != nil {
		panic(err)
	}

	var europeans = make([]*api_pb.European, len(res.Contracts))
	for index, contract := range res.Contracts {
		logrus.Info(contract.Strike)
		europeans[index] = &api_pb.European{
			Timestamp: res.Time[0].Timestamp,
			Undticker: "AAPL",
			Strike:    contract.Strike,
			Expiry:    contract.Calls[0].Expiry,
			Putcall:   contract.Putcall,
		}
	}

	return europeans
}
func insert(ctx context.Context, db *dgo.Dgraph, request interface{}, schema string) *api.Assigned {
	op := &api.Operation{
		Schema: schema,
	}
	err := db.Alter(ctx, op)
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
	assigned, err := db.NewTxn().Mutate(ctx, mu)
	if err != nil {
		logrus.Errorln(err)
		panic(err)
	}

	return assigned
}
