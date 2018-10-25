package cash

import (
	"fmt"

	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/x"
)

// ABCI Response Codes
// x/coins reserves 30 ~ 39.
const (
	CodeInsufficientFees  uint32 = 32
	CodeInsufficientFunds        = 33
	CodeInvalidAmount            = 34
	CodeInvalidMemo              = 35
	CodeEmptyAccount             = 36
)

var (
	errInsufficientFees  = fmt.Errorf("Insufficient fees")
	errInsufficientFunds = fmt.Errorf("Insufficient funds")
	errInvalidAmount     = fmt.Errorf("Invalid amount")
	errInvalidMemo       = fmt.Errorf("Invalid memo")
	errEmptyAccount      = fmt.Errorf("Account empty")
)

func ErrInsufficientFees(coin x.Coin) error {
	msg := coin.String()
	return errors.WithLog(msg, errInsufficientFees, CodeInsufficientFees)
}
func IsInsufficientFeesErr(err error) bool {
	return errors.IsSameError(errInsufficientFees, err)
}

func ErrInsufficientFunds() error {
	return errors.WithCode(errInsufficientFunds, CodeInsufficientFunds)
}
func IsInsufficientFundsErr(err error) bool {
	return errors.IsSameError(errInsufficientFunds, err)
}

func ErrInvalidAmount(reason string) error {
	return errors.WithLog(reason, errInvalidAmount, CodeInvalidAmount)
}
func IsInvalidAmountErr(err error) bool {
	return errors.IsSameError(errInvalidAmount, err)
}

func ErrInvalidMemo(reason string) error {
	return errors.WithLog(reason, errInvalidMemo, CodeInvalidMemo)
}
func IsInvalidMemoErr(err error) bool {
	return errors.IsSameError(errInvalidMemo, err)
}

func ErrEmptyAccount(addr weave.Address) error {
	return errors.WithLog(addr.String(), errEmptyAccount, CodeEmptyAccount)
}
func IsEmptyAccountErr(err error) bool {
	return errors.IsSameError(errEmptyAccount, err)
}
