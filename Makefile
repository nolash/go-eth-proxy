all: direct

direct:
	go build -v -o eth-proxy ./cmd/

geth:
	go build -v -o eth-proxy ./cmd/geth_proxy/

run: all
	./eth-proxy

clean:
	rm -vf eth-proxy

.PHONY: clean
