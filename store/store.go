package store

type Store interface {
	//GetTransaction(b []byte) (*types.Transaction, error)
	GetTransaction(b []byte) ([]byte, error)
	Close()
}

