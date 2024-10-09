# Name of the binary that will be created
BINARY_NAME = go-yaml-diff

# Source directory
SRC_DIR = .

# Build flags for the Go compiler
BUILD_FLAGS = -ldflags="-s -w"

# Default target to build and run the program
.PHONY: all
all: build run

# Compiles the program and creates the binary
.PHONY: build
build:
	@echo "Building the program..."
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(SRC_DIR)

# Runs the program
.PHONY: run
run: build
	@echo "Running the program..."
	./$(BINARY_NAME)

# Runs Go tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Cleans up the binary and other build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up build files..."
	rm -f $(BINARY_NAME)

# Tidies up the Go module dependencies
.PHONY: tidy
tidy:
	@echo "Tidying up Go module dependencies..."
	go mod tidy

# Installs the project dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download

# Updates dependencies to their latest versions
.PHONY: update
update:
	@echo "Updating dependencies to their latest versions..."
	go get -u ./...

# Help target to display available commands
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build     - Compiles the program and creates the binary"
	@echo "  make run       - Runs the program"
	@echo "  make test      - Runs Go tests"
	@echo "  make clean     - Removes build files and binaries"
	@echo "  make tidy      - Tidies up Go module dependencies"
	@echo "  make deps      - Downloads the project dependencies"
	@echo "  make update    - Updates dependencies to the latest versions"
	@echo "  make help      - Displays this help message"
