#!/usr/bin/env bash
set -e

# Gitingo global installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/devxdh/gitingo/main/install.sh | bash

REPO="devxdh/gitingo"
TAG="v1.0.0"
APP_NAME="gitingo"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Error: Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

case "$OS" in
    linux)
        OS="linux"
        ;;
    darwin)
        OS="darwin"
        ;;
    mingw*|msys*|cygwin*)
        OS="windows"
        ;;
    *)
        echo "Error: Unsupported operating system: $OS"
        exit 1
        ;;
esac

BINARY_NAME="${APP_NAME}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi

INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

echo "=========================================="
echo " Installing ${APP_NAME} (${OS}/${ARCH})..."
echo "=========================================="

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TAG}/${BINARY_NAME}"

# If building from source locally because release asset isn't fetched
if command -v go >/dev/null 2>&1; then
    echo "Found Go installed. Compiling latest ${APP_NAME} binary locally..."
    go install github.com/${REPO}@latest
    echo "Successfully installed ${APP_NAME} via go install!"
    echo "Location: $(command -v gitingo || echo '$GOPATH/bin/gitingo')"
    exit 0
fi

echo "Downloading binary from ${DOWNLOAD_URL}..."
if curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/gitingo"; then
    chmod +x "$TMP_DIR/gitingo"
    mv "$TMP_DIR/gitingo" "$INSTALL_DIR/gitingo"
    echo "Successfully installed ${APP_NAME} to ${INSTALL_DIR}/gitingo"
    echo "Run 'gitingo --help' to get started!"
else
    echo "Error: Download failed. Please ensure Go is installed and run:"
    echo "  go install github.com/${REPO}@latest"
    exit 1
fi
