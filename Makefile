BIN := ibugtool
VERSION := $(shell git describe --tags --always --dirty)

.PHONY: all clean $(BIN)

all: $(BIN)

$(BIN):
	go build -ldflags='-s -w -X main.version=$(VERSION)'

clean:
	rm -f $(BIN)
