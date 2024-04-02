BIN := ibugtool

.PHONY: all clean $(BIN)

all: $(BIN)

$(BIN):
	go build -ldflags='-s -w' .

clean:
	rm -f $(BIN)
