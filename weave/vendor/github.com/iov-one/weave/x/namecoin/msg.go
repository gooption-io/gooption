package namecoin

import (
	"regexp"

	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/x"
)

// Ensure we implement the Msg interface
var _ weave.Msg = (*NewTokenMsg)(nil)

const (
	pathNewTokenMsg       = "namecoin/ticker"
	pathSetNameMsg        = "namecoin/set_name"
	setNameCost     int64 = 50
	newTokenCost    int64 = 100

	minSigFigs = 0
	maxSigFigs = 9
)

var (
	// IsTokenName limits the human-readable names of the tokens,
	// subset of ASCII to avoid unicode tricks.
	IsTokenName = regexp.MustCompile(`^[A-Za-z0-9 \-_:]{3,32}$`).MatchString
	// IsWalletName is allowed names to attach to a wallet address
	IsWalletName = regexp.MustCompile(`^[a-z0-9_]{4,20}$`).MatchString
)

// Path returns the routing path for this message
func (NewTokenMsg) Path() string {
	return pathNewTokenMsg
}

// Validate makes sure that this is sensible
func (t *NewTokenMsg) Validate() error {
	if !x.IsCC(t.Ticker) {
		return x.ErrInvalidCurrency(t.Ticker)
	}
	if !IsTokenName(t.Name) {
		return ErrInvalidTokenName(t.Name)
	}
	if t.SigFigs < minSigFigs || t.SigFigs > maxSigFigs {
		return ErrInvalidSigFigs(t.SigFigs)
	}
	return nil
}

// BuildTokenMsg is a compact constructor for *NewTokenMsg
func BuildTokenMsg(ticker, name string, sigFigs int32) *NewTokenMsg {
	return &NewTokenMsg{
		Ticker:  ticker,
		Name:    name,
		SigFigs: sigFigs,
	}
}

// Path returns the routing path for this message
func (SetWalletNameMsg) Path() string {
	return pathSetNameMsg
}

// Validate makes sure that this is sensible
func (s *SetWalletNameMsg) Validate() error {
	if len(s.Address) != weave.AddressLength {
		return errors.ErrUnrecognizedAddress(s.Address)
	}
	if !IsWalletName(s.Name) {
		return ErrInvalidWalletName(s.Name)
	}
	return nil
}

// BuildSetNameMsg is a compact constructor for *SetWalletNameMsg
func BuildSetNameMsg(addr weave.Address, name string) *SetWalletNameMsg {
	return &SetWalletNameMsg{
		Address: addr,
		Name:    name,
	}
}
