package proxy

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

type ProxyService struct {
}

func (p *ProxyService) GetTransaction(ctx context.Context, hsh string) (*types.Transaction, error) {
	return nil, nil
}

