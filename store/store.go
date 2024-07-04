package store

type Store interface {
	GetTransaction(b []byte) ([]byte, error)
	GetBlock(b []byte) ([]byte, error)
	GetBlockNumber(n []byte) ([]byte, error)
	GetTransactionReceipt(b []byte) ([]byte, error)
	Close()
}
