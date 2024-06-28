all:
	go build -v -o eth-proxy ./cmd/

run: all
	./eth-proxy
