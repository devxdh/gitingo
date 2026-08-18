#!/usr/bin/env bash
set -e

APP_NAME="gitingo"
VERSION="${VERSION:-v1.0.0}"
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="${PROJECT_ROOT}/dist"
PACKAGE="github.com/devxdh/gitingo"
LDFLAGS="-s -w -X ${PACKAGE}/cmd.Version=${VERSION}"

rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

echo "Building $APP_NAME $VERSION..."

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    BINARY_NAME="${APP_NAME}-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi

    echo "Compiling for $GOOS/$GOARCH -> $DIST_DIR/$BINARY_NAME"
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$LDFLAGS" -o "$DIST_DIR/$BINARY_NAME" main.go

    # Create compressed archive packages for distribution
    TMP_ARCHIVE_DIR=$(mktemp -d)
    if [ "$GOOS" = "windows" ]; then
        cp "$DIST_DIR/$BINARY_NAME" "$TMP_ARCHIVE_DIR/gitingo.exe"
        (cd "$TMP_ARCHIVE_DIR" && zip -q -r "$DIST_DIR/gitingo_${VERSION}_${GOOS}_${GOARCH}.zip" gitingo.exe)
    else
        cp "$DIST_DIR/$BINARY_NAME" "$TMP_ARCHIVE_DIR/gitingo"
        (cd "$TMP_ARCHIVE_DIR" && tar -czf "$DIST_DIR/gitingo_${VERSION}_${GOOS}_${GOARCH}.tar.gz" gitingo)
    fi
    rm -rf "$TMP_ARCHIVE_DIR"
done

echo "Generating SHA256 checksums..."
(cd "$DIST_DIR" && (sha256sum * > checksums.txt 2>/dev/null || shasum -a 256 * > checksums.txt))

echo "Build complete! Release artifacts located in $DIST_DIR:"
ls -lh "$DIST_DIR"

