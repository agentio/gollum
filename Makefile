bootstrap:	
	go run ./cmd/bootstrap xrpc
	go run ./cmd/bootstrap xcli

slink:
	go install ./cmd/slink

all:	bootstrap slink	
