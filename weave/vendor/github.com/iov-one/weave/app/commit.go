package app

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
)

// commitStore is an internal type to handle loading from a
// KVCommitStore, maintaining different CacheWraps for
// Deliver and Check, and returning useful state info.
type commitStore struct {
	committed weave.CommitKVStore
	deliver   weave.KVCacheWrap
	check     weave.KVCacheWrap
}

// newCommitStore loads the CommitKVStore from disk or panics
// Sets up the deliver and check caches
func newCommitStore(store weave.CommitKVStore) *commitStore {
	err := store.LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	return &commitStore{
		committed: store,
		deliver:   store.CacheWrap(),
		check:     store.CacheWrap(),
	}
}

// CommitInfo returns the current height and hash
func (cs *commitStore) CommitInfo() (version int64, hash []byte) {
	id := cs.committed.LatestVersion()
	return id.Version, id.Hash
}

// Commit will flush deliver to the underlying store and commit it
// to disk. It then regenerates new deliver/check caches
//
// TODO: this should probably be protected by a mutex....
// need to think what concurrency we expect
func (cs *commitStore) Commit() weave.CommitID {
	// flush deliver to store and discard check
	cs.deliver.Write()
	cs.check.Discard()

	// write the store to disk
	res := cs.committed.Commit()

	// set up new caches
	cs.deliver = cs.committed.CacheWrap()
	cs.check = cs.committed.CacheWrap()
	return res
}

//------- storing chainID ---------

// _wv: is a prefix for weave internal data
const chainIDKey = "_wv:chainID"

// loadChainID returns the chain id stored if any
func loadChainID(kv weave.KVStore) string {
	v := kv.Get([]byte(chainIDKey))
	return string(v)
}

// saveChainID stores a chain id in the kv store.
// Returns error if already set, or invalid name
func saveChainID(kv weave.KVStore, chainID string) error {
	if !weave.IsValidChainID(chainID) {
		return errors.ErrInvalidChainID(chainID)
	}
	k := []byte(chainIDKey)
	if kv.Has(k) {
		return errors.ErrModifyChainID()
	}
	kv.Set(k, []byte(chainID))
	return nil
}
