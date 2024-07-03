package geth

import (
	"context"
	
	"github.com/go-ethereum/ethereum/core/types"

	"defalsify.org/go-eth-proxy/store"
)


type ProxyService struct {
	store store.Store
}

func NewProxyService(store store.Store) (*ProxyService) {
	return &ProxyService{
		store: store,
	}
}

func (p *ProxyService) GetTransactionByHash(ctx context.Context, hsh string) (*types.Transaction, error) {
	tx := &types.Transaction{}

	err = tx.UnmarshalJSON(b)
	if err != nil {
		return nil, err
	}
	log.Printf("tx %s gasprice %u gas %u", tx.Type(), tx.GasPrice(), tx.Gas())
	return tx, err

	return tx, nil
}
