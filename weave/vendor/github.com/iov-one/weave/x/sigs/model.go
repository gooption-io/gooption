package sigs

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/crypto"
	"github.com/iov-one/weave/orm"
)

// BucketName is where we store the accounts
const BucketName = "sigs"

//---- UserData
// Model stores the persistent state and all domain logic
// associated with valid state and state transitions.

var _ orm.CloneableData = (*UserData)(nil)

// Validate requires that all coins are in alphabetical
func (u *UserData) Validate() error {
	seq := u.Sequence
	if seq < 0 {
		return ErrInvalidSequence("Seq(%d)", seq)
	}
	if seq > 0 && u.PubKey == nil {
		return ErrInvalidSequence("Seq(%d) needs PubKey", seq)
	}
	return nil
}

// Copy makes a new UserData with the same coins
func (u *UserData) Copy() orm.CloneableData {
	return &UserData{
		Sequence: u.Sequence,
		PubKey:   u.PubKey,
	}
}

// CheckAndIncrementSequence checks if the current Sequence
// matches the expected value.
// If so, it will increase the sequence by one and return nil
// If not, it will not change the sequence, but return an error
func (u *UserData) CheckAndIncrementSequence(check int64) error {
	if u.Sequence != check {
		return ErrInvalidSequence("Mismatch %d != %d", check, u.Sequence)
	}
	u.Sequence++
	return nil
}

// SetPubKey will try to set the PubKey or panic on an illegal operation.
// It is illegal to reset an already set key
// Otherwise, we don't control
// (although we could verify the hash, we leave that to the controller)
func (u *UserData) SetPubKey(pubKey *crypto.PublicKey) {
	if u.PubKey != nil {
		panic("Cannot change pubkey for a user")
	}
	u.PubKey = pubKey
}

//-------------------- Object Wrapper -------

// AsUser will safely type-cast any value from Bucket to a UserData
func AsUser(obj orm.Object) *UserData {
	if obj == nil || obj.Value() == nil {
		return nil
	}
	return obj.Value().(*UserData)
}

// NewUser constructs an object from an address and pubkey
func NewUser(pubKey *crypto.PublicKey) orm.Object {
	var key weave.Address
	value := &UserData{PubKey: pubKey}
	if pubKey != nil {
		key = pubKey.Address()
	}
	return orm.NewSimpleObj(key, value)
}

//------------------ High-Level ------------------------

// Bucket extends orm.Bucket with GetOrCreate
type Bucket struct {
	orm.Bucket
}

// NewBucket creates the proper bucket for this extension
func NewBucket() Bucket {
	return Bucket{
		Bucket: orm.NewBucket(BucketName, NewUser(nil)),
	}
}

// GetOrCreate initializes a UserData if none exist for that key
func (b Bucket) GetOrCreate(db weave.KVStore,
	pubKey *crypto.PublicKey) (orm.Object, error) {

	obj, err := b.Get(db, pubKey.Address())
	if err == nil && obj == nil {
		obj = NewUser(pubKey)
	}
	return obj, err
}
