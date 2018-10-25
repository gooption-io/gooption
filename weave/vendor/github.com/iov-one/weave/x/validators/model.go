package validators

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/orm"
)

const (
	// BucketName contains address that are allowed to update validators
	BucketName = "uvalid"
	// Key is used to store account data
	Key = "accounts"
)

// WeaveAccounts is used to parse the json from genesis file
// use weave.Address, so address in hex, not base64
type WeaveAccounts struct {
	Addresses []weave.Address `json:"addresses"`
}

func (wa WeaveAccounts) Validate() error {
	for _, v := range wa.Addresses {
		err := v.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func AsWeaveAccounts(a *Accounts) WeaveAccounts {
	addrs := make([]weave.Address, len(a.Addresses))
	for k, v := range a.Addresses {
		addrs[k] = weave.Address(v)
	}
	return WeaveAccounts{Addresses: addrs}
}

func AsAccounts(a WeaveAccounts) *Accounts {
	addrs := make([][]byte, len(a.Addresses))
	for k, v := range a.Addresses {
		addrs[k] = []byte(v)
	}
	return &Accounts{Addresses: addrs}
}

// Copy makes new accounts object with the same addresses
func (m *Accounts) Copy() orm.CloneableData {
	addrSlice := make([][]byte, len(m.Addresses))
	for k, v := range m.Addresses {
		addr := make([]byte, len(v))
		copy(addr, v)
		addrSlice[k] = addr
	}
	return &Accounts{
		Addresses: addrSlice,
	}
}

func (m *Accounts) Validate() error {
	return AsWeaveAccounts(m).Validate()
}

func GetAccounts(bucket orm.Bucket, kv weave.KVStore) (*Accounts, error) {
	res, err := bucket.Get(kv, []byte(Key))
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, ErrWrongType(nil)
	}
	switch t := res.Value().(type) {
	case *Accounts:
		return t, nil
	default:
		return nil, ErrWrongType(t)
	}

}

func HasPermission(accts WeaveAccounts, checkAddress CheckAddress) bool {

	for _, v := range accts.Addresses {
		if checkAddress(v) {
			return true
		}
	}

	return false
}

func NewBucket() orm.Bucket {
	return orm.NewBucket(BucketName, NewAccounts())
}

func AccountsWith(acct WeaveAccounts) orm.Object {
	acc := AsAccounts(acct)
	return orm.NewSimpleObj([]byte(Key), acc)
}

// NewWallet creates an empty wallet with this address
// serves as an object for the bucket
func NewAccounts() orm.Object {
	return orm.NewSimpleObj([]byte(Key), new(Accounts))
}
