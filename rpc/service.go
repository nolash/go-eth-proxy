package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"defalsify.org/go-eth-proxy/store"
)

type LiteralProxyService struct {
	store store.Store
}

func NewProxyService(store store.Store) (*LiteralProxyService) {
	return &LiteralProxyService{
		store: store,
	}
}


func (p *LiteralProxyService) GetTransactionByHash(ctx context.Context, hsh string) ([]byte, error) {
	var err error

	b := []byte(hsh)
	b, err = p.store.GetTransaction(b)
	if err != nil {
		return nil, err
	}

	log.Printf("gettin for %x %s", hsh, b)
	j := &jsonRpcResponseFull{
		Jsonrpc: "2.0",
		Id: "",
		Result: string(b),
	}

	b = []byte{}
	r := bytes.NewBuffer(b)
	o := json.NewEncoder(r)
	o.Encode(j)

	return r.Bytes(), nil	
}
