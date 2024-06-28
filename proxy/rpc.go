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
	var uri *url.URL
	var err error

	if remoteURI != "" {
		uri, err = url.Parse(remoteURI)
		if err != nil {
			return nil, err
		}
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
	c, err := io.ReadFull(r.Body, b)
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
		"eth_getTransactionByHash",
	}) {
		if msg.Method == k {
			log.Printf("match %s", msg.Method)
			s.Server.ServeHTTP(w, r)
			return
		}
	}

	if s.uri == nil {
		log.Printf("missing remote side for unproxied method: %s", msg.Method)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	client_req := &http.Request{}
	client_req.Method = "POST"
	client_req.URL = s.uri
	client_req.Body = r.Body
	client_req.ContentLength = int64(c)
	client := &http.Client{}
	res, err := client.Do(client_req)
	if err != nil {
		log.Printf("%s", err)
		r.Body.Close()
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if res.StatusCode != http.StatusOK {
		v, _ := io.ReadAll(res.Body)
		log.Printf("%s", v)
	}
	w.WriteHeader(res.StatusCode)
	rr = io.TeeReader(res.Body, w)
	io.ReadAll(rr)
}
