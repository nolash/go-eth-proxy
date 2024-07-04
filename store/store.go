package store

type Store interface {
	GetTransaction(b []byte) ([]byte, error)
	Close()
}

