package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/rpc"
)

type jsonRpcMsg struct {
	Jsonrpc string
	Id string
	Method string
	Params []any

}
type ProxyServer struct {
	*rpc.Server
	uri *url.URL
}

func NewProxyServer(svc *ProxyService, remoteURI string) (*ProxyServer, error) {
	uri, err := url.Parse(remoteURI)
	if err != nil {
		return nil, err
	}
	srv := &ProxyServer{
		Server: rpc.NewServer(),
		uri: uri,
	}
	err = srv.Server.RegisterName("eth", svc)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func (s *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var rr io.Reader
	msg := jsonRpcMsg{}
	b := make([]byte, r.ContentLength)
	_, err := io.ReadFull(r.Body, b)
	if (err != nil) {
		log.Printf("%s", err)
		r.Body.Close()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	r.Body.Close()
	err = json.Unmarshal(b, &msg)
	if (err != nil) {
		log.Printf("%s", err)
		r.Body.Close()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rr = bytes.NewReader(b)
	r.Body = io.NopCloser(rr)

	for _, k := range([]string{
		"eth_getTransactionFromHash",
	}) {
		if msg.Method == k {
			log.Printf("match %s", msg.Method)
			s.Server.ServeHTTP(w, r)
			return
		}
	}
	client_req := r.Clone(r.Context())
	client_req.RequestURI = ""
	client_req.Method = "POST"
	client_req.URL = s.uri
	client := &http.Client{}
	res, err := client.Do(client_req)
	if err != nil {
		log.Printf("%s", err)
		r.Body.Close()
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	w.WriteHeader(res.StatusCode)
	rr = io.TeeReader(r.Body, w)
	io.ReadAll(rr)
}
