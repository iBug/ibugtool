BIN := ibugtool
VERSION := $(shell git describe --tags --always --dirty)

.PHONY: all gz clean $(BIN)

all: $(BIN)

gz: $(BIN).gz

$(BIN):
	go build -ldflags='-s -w -X main.version=$(VERSION)'

$(BIN).gz: $(BIN)
	gzip -c9 $^ > $@

clean:
	rm -f $(BIN)
