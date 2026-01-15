bootstrap:
	go install ./cmd/bootstrap
	go run ./cmd/bootstrap

slink:
	go install ./...

all:
	go run ./cmd/bootstrap
	go install ./...
