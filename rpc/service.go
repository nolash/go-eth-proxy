package rpc

import (
	"context"

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


func (p *LiteralProxyService) GetTransactionByHash(ctx context.Context, hsh string) ([]byte, error) {
	b := common.FromHex(hsh)
	tx, err := p.store.GetTransaction(b)
	if err != nil {
		return nil, err
	}

	b, err = tx.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return b, nil	
}
