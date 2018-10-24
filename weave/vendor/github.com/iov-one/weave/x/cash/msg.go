package cash

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/x"
)

// Ensure we implement the Msg interface
var _ weave.Msg = (*SendMsg)(nil)

const (
	pathSendMsg       = "cash/send"
	sendTxCost  int64 = 100

	maxMemoSize int = 128
	maxRefSize  int = 64
)

// Path returns the routing path for this message
func (SendMsg) Path() string {
	return pathSendMsg
}

// Validate makes sure that this is sensible
func (s *SendMsg) Validate() error {
	amt := s.GetAmount()
	if x.IsEmpty(amt) || !amt.IsPositive() {
		return ErrInvalidAmount("Non-positive SendMsg")
	}
	if err := amt.Validate(); err != nil {
		return err
	}
	if err := weave.Address(s.Src).Validate(); err != nil {
		return err
	}
	if err := weave.Address(s.Dest).Validate(); err != nil {
		return err
	}
	if len(s.GetMemo()) > maxMemoSize {
		return ErrInvalidMemo("Memo too long")
	}
	if len(s.GetRef()) > maxRefSize {
		return ErrInvalidMemo("Ref too long")
	}
	return nil
}

// DefaultSource makes sure there is a payer.
// If it was already set, returns s.
// If none was set, returns a new SendMsg with the source set
func (s *SendMsg) DefaultSource(addr []byte) *SendMsg {
	if len(s.GetSrc()) != 0 {
		return s
	}
	return &SendMsg{
		Src:    addr,
		Dest:   s.GetDest(),
		Amount: s.GetAmount(),
		Memo:   s.GetMemo(),
		Ref:    s.GetRef(),
	}
}

// FeeTx exposes information about the fees that
// should be paid
type FeeTx interface {
	GetFees() *FeeInfo
}

// DefaultPayer makes sure there is a payer.
// If it was already set, returns f.
// If none was set, returns a new FeeInfo, with the
// New address set
func (f *FeeInfo) DefaultPayer(addr []byte) *FeeInfo {
	if len(f.GetPayer()) != 0 {
		return f
	}
	return &FeeInfo{
		Payer: addr,
		Fees:  f.GetFees(),
	}
}

// Validate makes sure that this is sensible.
// Note that fee must be present, even if 0
func (f *FeeInfo) Validate() error {
	if f == nil {
		return errors.ErrUnrecognizedAddress(nil)
	}
	fee := f.GetFees()
	if fee == nil {
		return ErrInvalidAmount("Fees nil")
	}
	if err := fee.Validate(); err != nil {
		return err
	}
	if !fee.IsNonNegative() {
		return ErrInvalidAmount("Negative fees")
	}
	return weave.Address(f.Payer).Validate()
}
