package rpc


import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
//	"strconv"
)

type jsonRpcMsg struct {
	Method string

}

type jsonRpcError struct {
	Code int
}

type jsonRpcResponse struct {
	Error jsonRpcError
}

type ProxyServer struct {
	Server http.Handler
	uri *url.URL
}

type proxyWriter struct {
	header map[string][]string
	status int
	data *bytes.Buffer
	afterHeader bool
}


func (p *proxyWriter) Header() http.Header {
	return p.header
}

func (p *proxyWriter) Write(b []byte) (int, error) {
	log.Printf("proxyserver %s", b)
	return p.data.Write(b)
}

func (p *proxyWriter) WriteHeader(status int) {
	p.status = status

}

func (p *proxyWriter) Copy(w http.ResponseWriter) (int, error) {
	c := 0
	l := p.data.Len()
	b := p.data.Bytes()
	for ;c < l; {
		r, err := w.Write(b[c:])
		if err != nil {
			return 0, err
		}
		c += r
	}
	return c, nil
}

func newProxyWriter() *proxyWriter {
	b := make([]byte, 0, 1024)
	p := &proxyWriter{
		header: make(map[string][]string),
		data: bytes.NewBuffer(b),
	}
	return p
}

func NewProxyServer(backend http.Handler, remoteURI string) (*ProxyServer, error) {
	var uri *url.URL
	var err error

	if remoteURI != "" {
		uri, err = url.Parse(remoteURI)
		if err != nil {
			return nil, err
		}
	}
	srv := &ProxyServer{
		Server: backend,
		uri: uri,
	}
	return srv, nil
}

func (s *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//var rrr io.Reader
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

	rr := bytes.NewReader(b)
	r.Body = io.NopCloser(rr)

	for _, k := range([]string{
		"eth_getTransactionByHash",
	}) {
		rw := newProxyWriter()
		if msg.Method == k {
			s.Server.ServeHTTP(rw, r)
			rsp := jsonRpcResponse{}
			err = json.Unmarshal(b, &rsp)
			if (err != nil) {
				log.Printf("%s", err)
				r.Body.Close()
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//hd := rw.Header()
			//statusHd := hd["status"]
			//log.Printf("got status %s from proxy", statusHd)
			//if len(statusHd) > 0 && statusHd[0][:1] == "2" {
			if rsp.Error.Code == 0 {
//				statusCode, err := strconv.Atoi(hd["Status"][0])
//				if err != nil {
//					r.Body.Close()
//					w.WriteHeader(http.StatusInternalServerError)
//					return
//				}
//				rw.WriteHeader(statusCode)
				rw.WriteHeader(http.StatusOK)
				rw.Copy(w)
				return
			}

			log.Printf("not found in proxy: %s", k)
			rr.Seek(0, io.SeekStart)
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
	rrr := io.TeeReader(res.Body, w)
	io.ReadAll(rrr)
}
