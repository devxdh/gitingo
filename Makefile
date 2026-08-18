APP_NAME := gitingo
VERSION ?= v1.0.0
BUILD_DIR := bin
DIST_DIR := dist
PACKAGE := github.com/devxdh/gitingo
LDFLAGS := -s -w -X $(PACKAGE)/cmd.Version=$(VERSION)

.PHONY: all dev build install test clean cross-build help

all: test build

dev:
	go run main.go

build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) main.go
	@echo "Built $(APP_NAME) binary at $(BUILD_DIR)/$(APP_NAME)"

install: build
	@if [ -n "$$GOPATH" ]; then \
		cp $(BUILD_DIR)/$(APP_NAME) $$GOPATH/bin/$(APP_NAME); \
		echo "Installed $(APP_NAME) to $$GOPATH/bin/$(APP_NAME)"; \
	else \
		cp $(BUILD_DIR)/$(APP_NAME) $(HOME)/go/bin/$(APP_NAME) 2>/dev/null || cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/$(APP_NAME); \
		echo "Installed $(APP_NAME) globally!"; \
	fi

test:
	go test -v -cover ./...

cross-build:
	./build.sh

clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "Cleaned build artifacts."

help:
	@echo "Available make targets:"
	@echo "  make build       - Build local binary for current OS/arch"
	@echo "  make install     - Build and install binary globally"
	@echo "  make test        - Run unit tests with coverage"
	@echo "  make cross-build - Build release binaries for Linux, macOS, and Windows"
	@echo "  make dev         - Run application directly with go run"
	@echo "  make clean       - Remove binary and dist directories"
