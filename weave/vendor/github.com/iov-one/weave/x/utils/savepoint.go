package utils

import (
	"github.com/iov-one/weave"
)

// Savepoint will isolate all data inside of the call,
// and commit/rollback to savepoint based on if error
type Savepoint struct {
	onCheck   bool
	onDeliver bool
}

var _ weave.Decorator = Savepoint{}

// NewSavepoint creates a Savepoint decorator,
// but you must call OnCheck/OnDeliver so it will be triggered
func NewSavepoint() Savepoint {
	return Savepoint{}
}

// OnCheck returns a savepoint that will trigger on CheckTx
func (s Savepoint) OnCheck() Savepoint {
	return Savepoint{
		onCheck:   true,
		onDeliver: s.onDeliver,
	}
}

// OnDeliver returns a savepoint that will trigger on DeliverTx
func (s Savepoint) OnDeliver() Savepoint {
	return Savepoint{
		onCheck:   s.onCheck,
		onDeliver: true,
	}
}

// Check will optionally set a checkpoint
func (s Savepoint) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx,
	next weave.Checker) (weave.CheckResult, error) {

	if !s.onCheck {
		return next.Check(ctx, store, tx)
	}

	cstore, ok := store.(weave.CacheableKVStore)
	if !ok {
		return next.Check(ctx, store, tx)
	}

	cache := cstore.CacheWrap()
	res, err := next.Check(ctx, cache, tx)
	if err == nil {
		cache.Write()
	} else {
		cache.Discard()
	}
	return res, err
}

// Deliver will optionally set a checkpoint
func (s Savepoint) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx,
	next weave.Deliverer) (weave.DeliverResult, error) {

	if !s.onDeliver {
		return next.Deliver(ctx, store, tx)
	}

	cstore, ok := store.(weave.CacheableKVStore)
	if !ok {
		return next.Deliver(ctx, store, tx)
	}

	cache := cstore.CacheWrap()
	res, err := next.Deliver(ctx, cache, tx)
	if err == nil {
		cache.Write()
	} else {
		cache.Discard()
	}
	return res, err
}
