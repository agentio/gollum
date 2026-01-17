all:	bootstrap slink

bootstrap:	
	go run ./cmd/bootstrap lint
	go run ./cmd/bootstrap xrpc
	go run ./cmd/bootstrap xcli

slink:
	go install ./cmd/slink
