.PHONY: build
build:
	go build -o bin/brutus cmd/brutus/main.go

.PHONY: run
run:
	go run cmd/brutus/main.go

.PHONY: test
test:
	go test pkg/tokenizer/tokenizer_test.go
