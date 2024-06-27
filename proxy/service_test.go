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
	tx_test := "0x1c7770f04251de106344bc5e4c25a27143db9d40504045039d3d25c5b20b7740"

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
	err = client.Call(&tx, "eth_getTransaction", tx_test)
	if err != nil {
		t.Fatal(err)
	}
}
