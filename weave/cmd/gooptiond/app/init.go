package app

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/log"

	"github.com/iov-one/weave"
	"github.com/iov-one/weave/crypto"
	"github.com/iov-one/weave/x"
	"github.com/iov-one/weave/x/namecoin"
)

// GenInitOptions will produce some basic options for one rich
// account, to use for dev mode
//
// You can set
func GenInitOptions(args []string) (json.RawMessage, error) {
	ticker := "IOV"
	if len(args) > 0 {
		ticker = args[0]
		if !x.IsCC(ticker) {
			return nil, fmt.Errorf("Invalid ticker %s", ticker)
		}
	}

	var addr string
	if len(args) > 1 {
		addr = args[1]
	} else {
		// if no address provided, auto-generate one
		// and print out a recovery phrase
		bz, phrase, err := GenerateCoinKey()
		if err != nil {
			return nil, err
		}
		addr = hex.EncodeToString(bz)
		fmt.Println(phrase)
	}

	opts := fmt.Sprintf(`{
    "wallets": [
      {
        "address": "%s",
        "name": "admin",
        "coins": [
          {
            "whole": 123456789,
            "ticker": "%s"
          }
        ]
      }
    ],
    "tokens": [
      {
        "ticker": "%s",
        "name": "Main token of this chain",
        "sig_figs": 6
      }
    ]
  }`, addr, ticker, ticker)
	return []byte(opts), nil
}

// GenerateApp is used to create a stub for server/start.go command
func GenerateApp(home string, logger log.Logger, debug bool) (abci.Application, error) {
	// db goes in a subdir, but "" -> "" for memdb
	var dbPath string
	if home != "" {
		dbPath = filepath.Join(home, "bov.db")
	}

	// TODO: anyone can make a token????
	stack := Stack(x.Coin{}, nil)
	app, err := Application("mycoin", stack, TxDecoder, dbPath, debug)
	if err != nil {
		return nil, err
	}
	app.WithInit(namecoin.Initializer{})

	// set the logger and return
	app.WithLogger(logger)
	return app, nil
}

type output struct {
	Pubkey *crypto.PublicKey  `json:"pub_key"`
	Secret *crypto.PrivateKey `json:"secret"`
}

// GenerateCoinKey returns the address of a public key,
// along with a json representation of the keys.
// You can give coins to this address and
// import the keys in the js client to use them
func GenerateCoinKey() (weave.Address, string, error) {
	// XXX: we need to generate BIP39 recovery phrases in crypto
	privKey := crypto.GenPrivKeyEd25519()
	pubKey := privKey.PublicKey()
	addr := pubKey.Address()

	out := output{Pubkey: pubKey, Secret: privKey}
	keys, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, "", err
	}

	return addr, string(keys), nil
}
