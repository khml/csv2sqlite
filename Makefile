# Common settings
BINARY_NAME=csv2sqlite
BUILD_DIR=./bin
SRC_DIR=./src
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE=$(shell date -u '+%Y-%m-%d %H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X 'main.BuildDate=$(BUILD_DATE)'"

# Default platform settings (for local development)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

# Get the local platform-specific binary name
ifeq ($(GOOS),windows)
    LOCAL_BINARY=$(BINARY_NAME).exe
else
    LOCAL_BINARY=$(BINARY_NAME)
endif

# Default target
.PHONY: all
all: build

# Create build directory if it doesn't exist
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Build for the current platform (for development)
.PHONY: build
build: $(BUILD_DIR)
	cd $(SRC_DIR) && go build $(LDFLAGS) -o ../$(BUILD_DIR)/$(LOCAL_BINARY) -v && cd -
	@echo "Binary built at $(BUILD_DIR)/$(LOCAL_BINARY)"

# Run the application with sample data
.PHONY: run
run: build
	$(BUILD_DIR)/$(LOCAL_BINARY) read --csv $(SRC_DIR)/sample.csv --table sample --db sample.sqlite

# Format all Go code
.PHONY: fmt
fmt:
	cd $(SRC_DIR) && go fmt ./... && cd -

# Vet all Go code
.PHONY: vet
vet:
	cd $(SRC_DIR) && go vet ./... && cd -

# Clean up binaries and temporary files
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f *.zip
	rm -f $(BINARY_NAME)_*
	rm -f $(BINARY_NAME).exe
	rm -f sample.sqlite

# Install the binary to $GOPATH/bin
.PHONY: install
install: build
	cp $(BUILD_DIR)/$(LOCAL_BINARY) $(GOPATH)/bin/

# Build for Linux
.PHONY: build-linux
build-linux: $(BUILD_DIR)
	cd $(SRC_DIR) && GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ../$(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 -v && cd -
	cp $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 ./$(BINARY_NAME)_linux_amd64

# Build for macOS
.PHONY: build-darwin
build-darwin: $(BUILD_DIR)
	cd $(SRC_DIR) && GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o ../$(BUILD_DIR)/$(BINARY_NAME)_darwin_arm64 -v && cd -
	cp $(BUILD_DIR)/$(BINARY_NAME)_darwin_arm64 ./$(BINARY_NAME)_darwin_arm64

# Build for Windows
.PHONY: build-windows
build-windows: $(BUILD_DIR)
	cd $(SRC_DIR) && GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ../$(BUILD_DIR)/$(BINARY_NAME).exe -v && cd -
	cp $(BUILD_DIR)/$(BINARY_NAME).exe ./$(BINARY_NAME).exe

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows

# Package all builds
.PHONY: package
package: build-all
	zip $(BINARY_NAME).linux_amd64.zip $(BINARY_NAME)_linux_amd64
	zip $(BINARY_NAME).darwin_arm64.zip $(BINARY_NAME)_darwin_arm64
	zip $(BINARY_NAME).windows_amd64.zip $(BINARY_NAME).exe

# Build all artifacts for release
.PHONY: build-release-artifacts
build-release-artifacts: package
	@echo "Release artifacts created:"
	@ls -la *.zip

# Print help information
.PHONY: help
help:
	@echo "Make targets:"
	@echo "  all               - Default target, builds the binary for your platform"
	@echo "  build             - Build binary for the current platform"
	@echo "  run               - Build and run with sample data"
	@echo "  test              - Run all tests"
	@echo "  test-coverage     - Run tests with coverage report"
	@echo "  fmt               - Format all Go code"
	@echo "  vet               - Run Go vet on all packages"
	@echo "  clean             - Remove build artifacts"
	@echo "  install           - Install binary to GOPATH/bin"
	@echo "  build-linux       - Build binary for Linux"
	@echo "  build-darwin      - Build binary for macOS"
	@echo "  build-windows     - Build binary for Windows"
	@echo "  build-all         - Build binaries for all platforms"
	@echo "  package           - Build and package all binaries"
	@echo "  build-release-artifacts - Build all release artifacts (packages all platforms)"
	@echo "  help              - Show this help message"
