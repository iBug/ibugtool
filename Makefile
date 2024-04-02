BIN := ibugtool
VERSION := $(shell git describe --tags --always --dirty)

.PHONY: all gz test clean $(BIN)

all: $(BIN)

gz: $(BIN).gz

$(BIN):
	go build -ldflags='-s -w -X main.version=$(VERSION)'

$(BIN).gz: $(BIN)
	gzip -c9 $^ > $@

test:
	go test -v ./...

clean:
	rm -f $(BIN)
