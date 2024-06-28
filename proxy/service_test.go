package proxy

import (
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/core/types"

	"defalsify.org/go-eth-proxy/store/lmdb"
)

func TestProxyServerStart(t *testing.T) {
	var err error
	var tx types.Transaction
	tx_test := "0x60891c813816bb378ee8af428c5eb53b0479c980307d265e4abe39b4efd02e1d"

	dbpath, dbenv := os.LookupEnv("TEST_LMDB_DIR")
	if !dbenv {
		dbpath = "."
	}
	log.Printf("dbpath %s", dbpath)

	db, err := lmdb.NewStore(dbpath)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	svc := NewProxyService(db)

	srv := rpc.NewServer()
	err = srv.RegisterName("eth", svc)
	if err != nil {
		t.Fatal(err)
	}
	client := rpc.DialInProc(srv)
	mods, err := client.SupportedModules()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("mods %s", mods)
	err = client.Call(&tx, "eth_getTransactionByHash", tx_test)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tx %v", tx_test)
}
