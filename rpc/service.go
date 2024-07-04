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

func wrapResult(id string, result []byte) []byte {
	s := fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"id\":\"%s\",\"result\":%s}", id, result)
	return []byte(s)
}

func (p *LiteralProxyService) GetTransactionByHash(ctx context.Context, id string, hsh string) ([]byte, error) {
	var err error

	b := common.FromHex(hsh)
	b, err = p.store.GetTransaction(b)
	if err != nil {
		return nil, err
	}

	r := wrapResult(id, b)
	return r, nil
}

func (p *LiteralProxyService) GetTransactionReceipt(ctx context.Context, id string, hsh string) ([]byte, error) {
	var err error

	b := common.FromHex(hsh)
	b, err = p.store.GetTransactionReceipt(b)
	if err != nil {
		return nil, err
	}

	r := wrapResult(id, b)
	return r, nil
}

func (p *LiteralProxyService) GetBlockByHash(ctx context.Context, id string, hsh string) ([]byte, error) {
	var err error

	b := common.FromHex(hsh)
	b, err = p.store.GetBlock(b)
	if err != nil {
		return nil, err
	}

	r := wrapResult(id, b)
	return r, nil
}

func (p *LiteralProxyService) GetBlockByNumber(ctx context.Context, id string, numhex string) ([]byte, error) {
	b := common.FromHex(numhex)
	b, err := p.store.GetBlockNumber(b)
	if err != nil {
		return nil, err
	}

	r := wrapResult(id, b)
	return r, nil
}
