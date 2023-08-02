BIN := ibugtool

.PHONY: $(BIN)

$(BIN):
	go build -ldflags='-s -w' .
