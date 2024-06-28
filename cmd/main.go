package main

import (
	"flag"
	"log"
	"os"
	"net/http"
	"strings"

	"defalsify.org/go-eth-proxy/proxy"
	"defalsify.org/go-eth-proxy/store/lmdb"

)


func main() {
	dbpath := flag.String("cachepath", ".", "Path to lmdb data")
	host := flag.String("host", "0.0.0.0", "Remote host")
	port := flag.String("port", "8545", "Remote path")
	flag.Parse()

	db, err := lmdb.NewStore(*dbpath)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	defer db.Close()

	svc := proxy.NewProxyService(db)
	h, err := proxy.NewProxyServer(svc, flag.Arg(0))
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	srv := &http.Server{
		Handler: h,
		Addr: strings.Join([]string{*host, *port}, ":"),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}
