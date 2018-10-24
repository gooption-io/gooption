package multisig

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/x"
)

// Decorator checks multisig contract if available
type Decorator struct {
	auth   x.Authenticator
	bucket ContractBucket
}

var _ weave.Decorator = Decorator{}

// NewDecorator returns a default multisig decorator
func NewDecorator(auth x.Authenticator) Decorator {
	return Decorator{auth, NewContractBucket()}
}

// Check enforce multisig contract before calling down the stack
func (d Decorator) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx, next weave.Checker) (weave.CheckResult, error) {
	var res weave.CheckResult
	newCtx, err := d.withMultisig(ctx, store, tx)
	if err != nil {
		return res, err
	}

	return next.Check(newCtx, store, tx)
}

// Deliver enforces multisig contract before calling down the stack
func (d Decorator) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx, next weave.Deliverer) (weave.DeliverResult, error) {
	var res weave.DeliverResult
	newCtx, err := d.withMultisig(ctx, store, tx)
	if err != nil {
		return res, err
	}

	return next.Deliver(newCtx, store, tx)
}

func (d Decorator) withMultisig(ctx weave.Context, store weave.KVStore, tx weave.Tx) (weave.Context, error) {
	if multisigContract, ok := tx.(MultiSigTx); ok {
		ids := multisigContract.GetMultisig()
		for _, contractID := range ids {
			if contractID == nil {
				return ctx, nil
			}

			// check if we already have it
			if d.auth.HasAddress(ctx, MultiSigCondition(contractID).Address()) {
				return ctx, nil
			}

			// load contract
			contract, err := d.getContract(store, contractID)
			if err != nil {
				return ctx, err
			}

			// collect all sigs
			sigs := make([]weave.Address, len(contract.Sigs))
			for i, sig := range contract.Sigs {
				sigs[i] = sig
			}

			// check sigs (can be sig or multisig)
			authenticated := x.HasNAddresses(ctx, d.auth, sigs, int(contract.ActivationThreshold))
			if !authenticated {
				return ctx, ErrUnauthorizedMultiSig(contractID)
			}

			ctx = withMultisig(ctx, contractID)
		}
	}

	return ctx, nil
}

func (d Decorator) getContract(store weave.KVStore, id []byte) (*Contract, error) {
	obj, err := d.bucket.Get(store, id)
	if err != nil {
		return nil, err
	}

	if obj == nil || (obj != nil && obj.Value() == nil) {
		return nil, ErrContractNotFound(id)
	}

	contract := obj.Value().(*Contract)
	return contract, err
}
