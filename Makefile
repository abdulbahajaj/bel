

.PHONY: build
build:
	go build -o bin/bel cmd/bel/main.go

.PHONY: run
run:
	go run cmd/bel/main.go
