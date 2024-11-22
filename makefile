# Variables
APP_NAME := mole
BUILD_DIR := build
OUTPUT_DIR := $(BUILD_DIR)/linux-amd64
OUT_DEV := $(BUILD_DIR)/dev
GOOS := linux
GOARCH := amd64

# Build the app for production
buildprod:
	@echo "Building $(APP_NAME) for production..."
	mkdir -p $(OUTPUT_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o $(OUTPUT_DIR)/$(APP_NAME) \
		-ldflags "-X 'github.com/zulubit/mole/pkg/consts.BasePath=/home/mole' -X 'github.com/zulubit/mole/pkg/consts.Prod=1'" ./cmd/cli

# Build the app for development
builddev:
	@echo "Building $(APP_NAME) for development..."
	mkdir -p $(OUTPUT_DIR)
	go build -o $(OUT_DEV)/$(APP_NAME) ./cmd/cli

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

