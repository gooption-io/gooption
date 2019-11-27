package server

import (
	"io/ioutil"
	"log"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/lehajam/dgo"
	"github.com/lehajam/dgo/protos/api"

	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

func insertObj(ctx context.Context, db *dgo.Dgraph, obj interface{}) (*api.Assigned, error) {
	txn := db.NewTxn()
	defer txn.Discard(context.Background())

	pb, err := json.Marshal(obj)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	assigned, err := txn.Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: pb})
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	return assigned, nil
}

func alterSchema(ctx context.Context, db *dgo.Dgraph, schema string) error {
	op := &api.Operation{
		Schema: schema,
	}

	return db.Alter(ctx, op)
}
