package proxy

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"

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
	log.Printf("get tx hash %s", hsh)
	b := common.FromHex(hsh)
	tx, err := p.store.GetTransaction(b)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
