package server

import (
	"io/ioutil"
	"log"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	dgo "github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"

	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// TODO: close connection
func newDgraphClient() *dgo.Dgraph {
	_, err := ioutil.TempDir("", "client_")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	dc := api.NewDgraphClient(conn)
	return dgo.NewDgraphClient(dc)
}

func insertObj(ctx context.Context, db *dgo.Dgraph, obj interface{}) (*api.Response, error) {
	txn := db.NewTxn()
	defer txn.Discard(context.Background())

	pb, err := json.Marshal(obj)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	resp, err := txn.Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: pb})
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	return resp, nil
}

func alterSchema(ctx context.Context, db *dgo.Dgraph, schema string) error {
	op := &api.Operation{
		Schema: schema,
	}

	return db.Alter(ctx, op)
}
