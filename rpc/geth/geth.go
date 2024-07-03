package geth

import (
	"net/http"

	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"defalsify.org/go-eth-proxy/store"
)

func NewGethBackend(db store.Store) (http.Handler, error) {
	var err error

	svc := NewProxyService(db)
	srv := gethrpc.NewServer()
	err = srv.RegisterName("eth", svc)
	if err != nil {
		return nil, err
	}
	return srv, nil

}
