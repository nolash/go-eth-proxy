package rpc

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

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


func (p *LiteralProxyService) GetTransactionByHash(ctx context.Context, id string, hsh string) ([]byte, error) {
	var err error

	b := common.FromHex(hsh)
	b, err = p.store.GetTransaction(b)
	if err != nil {
		return nil, err
	}

	s := fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"id\":\"%s\",\"result\":%s}", id, b)

	return []byte(s), nil
}
