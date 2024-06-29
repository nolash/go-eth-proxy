package lmdb

import (
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ledgerwatch/lmdb-go/lmdb"

	"defalsify.org/go-eth-proxy/store"
)

type LmdbStore struct {
	store store.Store
	env *lmdb.Env
	dbi lmdb.DBI
}

/// TODO: not create
func NewStore(path string) (*LmdbStore, error) {
	var err error

	log.Printf("lmdb store path: %s", path)
	o := &LmdbStore{}
	o.env, err = lmdb.NewEnv()
	if err != nil {
		return nil, err
	}
	err = o.env.SetMaxDBs(1)
	if err != nil {
		return nil, err
	}
	err = o.env.SetMapSize(1 << 30)
	if err != nil {
		return nil, err
	}
	err = o.env.Open(path, 0, 0644)
	if err != nil {
		return nil, err
	}
	err = o.env.Update(func(txn *lmdb.Txn) (error) {
		var err error
		o.dbi, err = txn.OpenRoot(0)
		return err
	})
	if err != nil {
		return nil, err
	}
	return o, nil
}


func (l *LmdbStore) GetTransaction(k []byte) (*types.Transaction, error) {
	var b []byte
	tx := &types.Transaction{}

	kp := make([]byte, len(k) + 7)
	copy(kp, []byte("tx/src/"))
	copy(kp[7:], k)

	err := l.env.View(func(txn *lmdb.Txn) (error) {
		log.Printf("gettx %x %v", kp, txn)
		v, err := txn.Get(l.dbi, kp)
		if err != nil {
			return err
		}
		b = make([]byte, len(v))
		copy(b, v)
		return nil
	})
	if err != nil {
		return tx, err
	}
	log.Printf("lmdn result: %s", b)
	err = tx.UnmarshalJSON(b)
	if err != nil {
		return nil, err
	}
	log.Printf("tx %s gasprice %u gas %u", tx.Type(), tx.GasPrice(), tx.Gas())
	return tx, err
}


func (l *LmdbStore) Close() {
	l.env.Close()
}
