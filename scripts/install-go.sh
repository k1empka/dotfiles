#!/bin/sh
set -e

# Install Go on macOS or Linux.
# Usage: sh scripts/install-go.sh

GO_VERSION="1.24.0"

if command -v go >/dev/null 2>&1; then
    echo "==> Go already installed: $(go version)"
    exit 0
fi

echo "==> Installing Go ${GO_VERSION}..."

case "$(uname -s)" in
    Darwin)
        if command -v brew >/dev/null 2>&1; then
            brew install go
        else
            echo "    Homebrew not found. Install it first: https://brew.sh"
            exit 1
        fi
        ;;
    Linux)
        ARCH="$(uname -m)"
        case "$ARCH" in
            x86_64)  GOARCH="amd64" ;;
            aarch64) GOARCH="arm64" ;;
            armv6l)  GOARCH="armv6l" ;;
            *)
                echo "    Unsupported architecture: $ARCH"
                exit 1
                ;;
        esac

        TARBALL="go${GO_VERSION}.linux-${GOARCH}.tar.gz"
        URL="https://go.dev/dl/${TARBALL}"

        echo "    Downloading ${URL}..."
        curl -fsSL -o "/tmp/${TARBALL}" "$URL"

        echo "    Installing to /usr/local/go..."
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf "/tmp/${TARBALL}"
        rm -f "/tmp/${TARBALL}"

        # Add to PATH for the current session.
        export PATH="/usr/local/go/bin:$PATH"

        echo "    Add the following to your shell profile:"
        echo "    export PATH=/usr/local/go/bin:\$PATH"
        ;;
    *)
        echo "    Unsupported OS: $(uname -s)"
        echo "    Download Go manually from https://go.dev/dl/"
        exit 1
        ;;
esac

echo "==> Installed: $(go version)"
