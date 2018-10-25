package validators

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/x"
)

type AuthCheckAddress = func(auth x.Authenticator, ctx weave.Context) CheckAddress

var authCheckAddress = func(auth x.Authenticator, ctx weave.Context) CheckAddress {
	return func(addr weave.Address) bool {
		return auth.HasAddress(ctx, addr)
	}
}

// RegisterRoutes will instantiate and register
// all handlers in this package
func RegisterRoutes(r weave.Registry, auth x.Authenticator,
	control Controller) {

	r.Handle(pathUpdate, NewUpdateHandler(auth, control, authCheckAddress))
}

// RegisterQuery will register this bucket as "/validators"
func RegisterQuery(qr weave.QueryRouter) {
	NewBucket().Register("validators", qr)
}

// UpdateHandler will handle sending coins
type UpdateHandler struct {
	auth             x.Authenticator
	control          Controller
	authCheckAddress AuthCheckAddress
}

var _ weave.Handler = UpdateHandler{}

// NewUpdateHandler creates a handler for SendMsg
func NewUpdateHandler(auth x.Authenticator, control Controller, checkAddr AuthCheckAddress) UpdateHandler {
	return UpdateHandler{
		auth:             auth,
		control:          control,
		authCheckAddress: checkAddr,
	}
}

// Check verifies all the preconditions
func (h UpdateHandler) Check(ctx weave.Context, store weave.KVStore,
	tx weave.Tx) (weave.CheckResult, error) {

	var res weave.CheckResult
	rmsg, err := tx.GetMsg()
	if err != nil {
		return res, err
	}
	msg, ok := rmsg.(*SetValidatorsMsg)
	if !ok {
		return res, errors.ErrUnknownTxType(rmsg)
	}

	_, err = h.control.CanUpdateValidators(store, h.authCheckAddress(h.auth, ctx), msg.AsABCI())
	if err != nil {
		return res, err
	}

	return res, nil
}

// Deliver provides the diff given everything is okay with permissions and such
// Check did the same job already, so we can assume stuff goes okay
func (h UpdateHandler) Deliver(ctx weave.Context, store weave.KVStore,
	tx weave.Tx) (weave.DeliverResult, error) {

	// ensure type and validate...
	var res weave.DeliverResult
	rmsg, err := tx.GetMsg()
	if err != nil {
		return res, err
	}
	msg, ok := rmsg.(*SetValidatorsMsg)
	if !ok {
		return res, errors.ErrUnknownTxType(rmsg)
	}

	diff, err := h.control.CanUpdateValidators(store, h.authCheckAddress(h.auth, ctx), msg.AsABCI())
	if err != nil {
		return res, err
	}

	res.Diff = diff

	return res, nil
}
