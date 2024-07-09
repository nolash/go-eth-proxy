package geth

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
	tx := &types.Transaction{}

	b := common.FromHex(hsh)
	r, err := p.store.GetTransaction(b)

	err = tx.UnmarshalJSON(r)
	if err != nil {
		return nil, err
	}
	log.Printf("tx %v gasprice %d gas %d", tx.Type(), tx.GasPrice(), tx.Gas())
	return tx, nil
}
