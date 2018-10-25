package gooption

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/orm"
	"github.com/iov-one/weave/x"
)

var (
	pathCreateValuationMsg       = "valuation/create"
	creationCost           int64 = 100
)

// RegisterRoutes will instantiate and register
// all handlers in this package
func RegisterRoutes(r weave.Registry, auth x.Authenticator) {
	bucket := NewValuationBucket()
	r.Handle(pathCreateValuationMsg, CreateValuationMsgHandler{auth, bucket})
}

// RegisterQuery register queries from buckets in this package
func RegisterQuery(qr weave.QueryRouter) {
	NewValuationBucket().Register("valuations", qr)
}

type CreateValuationMsgHandler struct {
	auth   x.Authenticator
	bucket ValuationBucket
}

var _ weave.Handler = CreateValuationMsgHandler{}

// Path returns the routing path for this message
func (CreateValuationMsg) Path() string {
	return pathCreateValuationMsg
}

func (h CreateValuationMsgHandler) Check(ctx weave.Context, db weave.KVStore, tx weave.Tx) (weave.CheckResult, error) {
	var res weave.CheckResult
	_, err := h.validate(ctx, db, tx)
	if err != nil {
		return res, err
	}

	res.GasAllocated = creationCost
	return res, nil
}

func (h CreateValuationMsgHandler) Deliver(ctx weave.Context, db weave.KVStore, tx weave.Tx) (weave.DeliverResult, error) {
	var res weave.DeliverResult
	msg, err := h.validate(ctx, db, tx)
	if err != nil {
		return res, err
	}

	valuation := &Valuation{
		Sender:     msg.Sender,
		Publisher:  msg.Publisher,
		ContractId: msg.ContractId,
		Timestamp:  msg.Timestamp,
		Request:    msg.Request,
		Response:   msg.Response,
	}

	id := h.bucket.idSeq.NextVal(db)
	obj := orm.NewSimpleObj(id, valuation)
	err = h.bucket.Save(db, obj)
	if err != nil {
		return res, err
	}

	res.Data = id
	return res, nil
}

// validate does all common pre-processing between Check and Deliver
func (h CreateValuationMsgHandler) validate(ctx weave.Context, db weave.KVStore, tx weave.Tx) (*CreateValuationMsg, error) {
	rmsg, err := tx.GetMsg()
	if err != nil {
		return nil, err
	}
	msg, ok := rmsg.(*CreateValuationMsg)
	if !ok {
		return nil, errors.ErrUnknownTxType(rmsg)
	}

	err = msg.Validate()
	if err != nil {
		return nil, err
	}

	// make sure we have permission from the sender
	if !h.auth.HasAddress(ctx, msg.Sender) {
		return nil, errors.ErrUnauthorized()
	}

	return msg, nil
}
