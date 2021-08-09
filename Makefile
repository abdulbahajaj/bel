.PHONY: build
build:
	go build -o bin/brut cmd/brutus/main.go

.PHONY: run
run:
	go run cmd/brutus/main.go code/core.brut

.PHONY: test
test:
	go run cmd/test/main.go
