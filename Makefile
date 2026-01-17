bootstrap:	
	go run ./cmd/bootstrap xrpc
	go run ./cmd/bootstrap slink

slink:
	go install ./cmd/slink

all:	bootstrap slink	
