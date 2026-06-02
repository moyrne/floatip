.PHONY: build run test clean

BIN_DIR=bin
BINARY=$(BIN_DIR)/floatip

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BINARY) ./cmd/floatip/

run:
	go run ./cmd/floatip/

test:
	go test ./...

clean:
	rm -rf $(BIN_DIR)
