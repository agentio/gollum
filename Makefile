all:
	go install ./cmd/generate-gollum
	generate-gollum
	go install ./...
