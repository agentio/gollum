bootstrap:
	go install ./cmd/bootstrap
	go run ./cmd/bootstrap
	go run ./cmd/bootstrap cli

slink:
	go install ./cmd/slink

all:	bootstrap slink
