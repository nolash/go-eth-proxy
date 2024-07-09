package rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type literalBackend struct {
	svc *LiteralProxyService
}

func NewBackend(svc *LiteralProxyService) *literalBackend {
	return &literalBackend {
		svc: svc,
	}

}

func inJson(b []byte) (*jsonRpcMsgFull, error) {
	msg := &jsonRpcMsgFull{}
	err := json.Unmarshal(b, msg)
	return msg, err
}

func (l *literalBackend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s := fmt.Sprintf("Status: %d Not valid jsonrpc request", http.StatusInternalServerError)
		w.Write([]byte(s))
		return
	}

	msg, err := inJson(b)
	if err != nil {
		s := fmt.Sprintf("Status: %d Not valid jsonrpc request", http.StatusBadRequest)
		w.Write([]byte(s))
		return
	}

	// TODO: make sure getblockbynumber is hex input
	switch msg.Method {
		case "eth_getTransactionByHash":
			b, err = l.svc.GetTransactionByHash(r.Context(), msg.Id, msg.Params[0].(string))
		case "eth_getTransactionReceipt":
			b, err = l.svc.GetTransactionReceipt(r.Context(), msg.Id, msg.Params[0].(string))
		case "eth_getBlockByHash":
			b, err = l.svc.GetBlockByHash(r.Context(), msg.Id, msg.Params[0].(string))
		case "eth_getBlockByNumber":
			b, err = l.svc.GetBlockByNumber(r.Context(), msg.Id, msg.Params[0].(string))
		default:
			s := fmt.Sprintf("Status: %d Method not supported", http.StatusBadRequest)
			w.Write([]byte(s))
			return
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(b))
}
