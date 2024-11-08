.PHONY: all build run test docker-build

BIN_DIR := bin
COVERAGE_DIR := coverage

all: build

build:
	@mkdir -p $(BIN_DIR)
	GOOS=linux go build -o $(BIN_DIR)/huproxy .

run: build
	@./$(BIN_DIR)/huproxy

test:
	go test -v ./...

coverage:
	@mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -func=$(COVERAGE_DIR)/coverage.out

coverage-html:
	@mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

clean:
	@rm -rf $(BIN_DIR) $(COVERAGE_DIR)
