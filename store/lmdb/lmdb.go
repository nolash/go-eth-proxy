package lmdb

import (
	"encoding/binary"
	"log"

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

func (l *LmdbStore) get(pfx string, k []byte) ([]byte, error) {
	var b []byte

	n := len(pfx)
	kp := make([]byte, len(k) + n)
	copy(kp, []byte(pfx))
	copy(kp[n:], k)

	err := l.env.View(func(txn *lmdb.Txn) (error) {
		log.Printf("get %x %v", kp, txn)
		v, err := txn.Get(l.dbi, kp)
		if err != nil {
			return err
		}
		b = make([]byte, len(v))
		copy(b, v)
		return nil
	})
	log.Printf("lmdb result: %s", b)
	if err != nil {
		return nil, err
	}
	return b, nil
}


func (l *LmdbStore) GetTransaction(k []byte) ([]byte, error) {
	return l.get("tx/src/", k)
}

func (l *LmdbStore) GetBlockNumber(n uint64) ([]byte, error) {
	var err error
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	
	b, err = l.get("block/num/", b)
	if err != nil {
		return nil, err
	}
	return l.get("block/hash/", b)
}

func (l *LmdbStore) GetBlockHash(k []byte) ([]byte, error) {
	return l.get("block/hash/", k)
}

func (l *LmdbStore) Close() {
	l.env.Close()
}
