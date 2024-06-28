package main

import (
	"flag"
	"log"
	"os"
	"net/http"

	"defalsify.org/go-eth-proxy/proxy"
	"defalsify.org/go-eth-proxy/store/lmdb"
	"github.com/ethereum/go-ethereum/rpc"

)


func main() {
	dbpath := flag.String("cachepath", ".", "Path to lmdb data")
	db, err := lmdb.NewStore(*dbpath)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer db.Close()

	svc := proxy.NewProxyService(db)
	h := rpc.NewServer()
	err = h.RegisterName("eth", svc)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Handler: h,
		Addr: "0.0.0.0:8080",
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}
