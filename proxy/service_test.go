package proxy

import (
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestProxyServerStart(t *testing.T) {
	var err error
	var tx types.Transaction

	srv := rpc.NewServer()
	svc := ProxyService{}
	err = srv.RegisterName("eth", &svc)
	if err != nil {
		t.Error(err)
	}
	client := rpc.DialInProc(srv)
	mods, err := client.SupportedModules()
	if err != nil {
		t.Error(err)
	}
	t.Logf("mods %s", mods)
	err = client.Call(&tx, "eth_getTransaction", "foo")
	if err != nil {
		t.Error(err)
	}
}
