package hashlock

import (
	"crypto/sha256"

	"github.com/iov-one/weave"
)

// HashKeyTx is an optional interface for a Tx that allows
// it to provide Keys (Preimages) to open HashLocks
type HashKeyTx interface {
	// GetPreimage should return a hash preimage if provided
	// or nil if not included in this tx
	GetPreimage() []byte
}

// PreimageCondition calculates a sha256 hash and then
func PreimageCondition(preimage []byte) weave.Condition {
	h := sha256.Sum256(preimage)
	return weave.NewCondition("hash", "sha256", h[:])
}
