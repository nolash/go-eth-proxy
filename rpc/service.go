package rpc

import (
	"context"
	"fmt"
	"log"

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

func wrapResult(id any, result []byte) ([]byte, error) {
	var s string
	var id_i int
	id_s, ok := id.(string)
	if ok {
		s = fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"id\":\"%s\",\"result\":%s}", id_s, result)
	} else {
		id_i, ok = id.(int)
		if !ok {
			return nil, fmt.Errorf("id not valid type")
		}
		s = fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"id\":%d,\"result\":%s}", id_i, result)
	}
	return []byte(s), nil
}

func (p *LiteralProxyService) GetTransactionByHash(ctx context.Context, id any, hsh string) ([]byte, error) {
	var err error

	b := common.FromHex(hsh)
	b, err = p.store.GetTransaction(b)
	if err != nil {
		return nil, err
	}

	return wrapResult(id, b)
}

func (p *LiteralProxyService) GetTransactionReceipt(ctx context.Context, id any, hsh string) ([]byte, error) {
	var err error

	b := common.FromHex(hsh)
	b, err = p.store.GetTransactionReceipt(b)
	if err != nil {
		return nil, err
	}

	return wrapResult(id, b)
}

func (p *LiteralProxyService) GetBlockByHash(ctx context.Context, id any, hsh string) ([]byte, error) {
	var err error

	b := common.FromHex(hsh)
	b, err = p.store.GetBlock(b)
	if err != nil {
		return nil, err
	}

	return wrapResult(id, b)
}

func (p *LiteralProxyService) GetBlockByNumber(ctx context.Context, id any, numhex string) ([]byte, error) {
	b := common.FromHex(numhex)
	b, err := p.store.GetBlockNumber(b)
	log.Printf("result %v", b)
	if err != nil {
		return nil, err
	}

	return wrapResult(id, b)
}
