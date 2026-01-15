bootstrap:
	go install ./cmd/bootstrap
	go run ./cmd/bootstrap

slink:
	go install ./cmd/slink

all:	bootstrap slink
