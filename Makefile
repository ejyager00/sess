BIN=sess
OUTPUT_DIR=target

build:
	go build -o $(OUTPUT_DIR)/$(BIN) main.go