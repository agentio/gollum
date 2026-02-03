all:	bootstrap slink

bootstrap:	
	go run ./cmd/bootstrap lint
	go run ./cmd/bootstrap xrpc
	go run ./cmd/bootstrap call
	go run ./cmd/bootstrap check

slink:
	go install -tags jwx_es256k ./cmd/slink

manifest:
	go run ./cmd/bootstrap lint
	go run ./cmd/bootstrap xrpc -m sample-manifest.json
	go run ./cmd/bootstrap call -m sample-manifest.json
	go run ./cmd/bootstrap check -m sample-manifest.json
	go install ./cmd/slink
