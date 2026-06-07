#!/bin/sh
set -e

REPO="ev-the-dev/gitgogit"
BINARY="gitgogit"
INSTALL_DIR="/usr/local/bin"

# Allow version to be overridden: VERSION=v1.0.0 ./install.sh
if [ -z "$VERSION" ]; then
    VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep '"tag_name"' \
        | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "error: could not determine latest release version" >&2
        exit 1
    fi
fi

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    linux)  OS="linux"  ;;
    darwin) OS="darwin" ;;
    *)
        echo "error: unsupported OS: $OS" >&2
        exit 1
        ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)          ARCH="amd64" ;;
    aarch64 | arm64) ARCH="arm64" ;;
    *)
        echo "error: unsupported architecture: $ARCH" >&2
        exit 1
        ;;
esac

ASSET="${BINARY}-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"

echo "Installing ${BINARY} ${VERSION} (${OS}/${ARCH})..."

curl -fsSL "$URL" -o "/tmp/${BINARY}"
chmod +x "/tmp/${BINARY}"

if [ -w "$INSTALL_DIR" ]; then
    mv "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
    echo "Moving binary to ${INSTALL_DIR} (may require sudo)..."
    sudo mv "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi

echo "Installed to ${INSTALL_DIR}/${BINARY}"
echo "Run '${BINARY} --version' to verify."
