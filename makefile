# Variables
APP_NAME := mole
BUILD_DIR := build
OUTPUT_DIR := $(BUILD_DIR)/linux
GOOS := linux
GOARCH := amd64

# Build the app for Linux with production environment
buildprod:
	@echo "Building $(APP_NAME) for production..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o $(OUTPUT_DIR)/$(APP_NAME) -ldflags "-X github.com/zulubit/mole/pkg/consts.BasePath=/home/mole" ./cmd/cli