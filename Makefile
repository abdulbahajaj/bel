.PHONY: build
build:
	go build -o bin/brut cmd/brutus/main.go

.PHONY: run
run:
	go run cmd/brutus/main.go tests/manual-test.brut

.PHONY: test
test:
	go test pkg/tokenizer/tokenizer_test.go
