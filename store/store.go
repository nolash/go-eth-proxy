package store

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type Store interface {
	GetTransaction(b []byte) (*types.Transaction, error)
	Close()
}
