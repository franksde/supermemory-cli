#!/bin/bash
set -euo pipefail

REPO="franksde/supermemory-cli"
BINARY="sm"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  *)       echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
  darwin|linux) ;;
  *)            echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Get latest release tag
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$LATEST" ]; then
  echo "Failed to fetch latest release"
  exit 1
fi

echo "Installing ${BINARY} ${LATEST} for ${OS}/${ARCH}..."

# Download
TARBALL="${BINARY}_${LATEST#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${TARBALL}"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

curl -fsSL "$URL" -o "${TMPDIR}/${TARBALL}"
curl -fsSL "https://github.com/${REPO}/releases/download/${LATEST}/checksums.txt" -o "${TMPDIR}/checksums.txt"

EXPECTED=$(awk -v file="$TARBALL" '$2 == file {print $1}' "${TMPDIR}/checksums.txt")
if [ -z "$EXPECTED" ]; then
  echo "Checksum for ${TARBALL} not found"
  exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
  ACTUAL=$(sha256sum "${TMPDIR}/${TARBALL}" | awk '{print $1}')
elif command -v shasum >/dev/null 2>&1; then
  ACTUAL=$(shasum -a 256 "${TMPDIR}/${TARBALL}" | awk '{print $1}')
else
  echo "No SHA-256 checksum tool found (need sha256sum or shasum)"
  exit 1
fi

if [ "$ACTUAL" != "$EXPECTED" ]; then
  echo "Checksum verification failed for ${TARBALL}"
  exit 1
fi

tar xzf "${TMPDIR}/${TARBALL}" -C "$TMPDIR"

# Install
if [ ! -d "$INSTALL_DIR" ]; then
  if ! mkdir -p "$INSTALL_DIR" 2>/dev/null; then
    echo "Need sudo to create ${INSTALL_DIR}"
    sudo mkdir -p "$INSTALL_DIR"
  fi
fi

if [ -w "$INSTALL_DIR" ]; then
  mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
  echo "Need sudo to install to ${INSTALL_DIR}"
  sudo mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi

chmod +x "${INSTALL_DIR}/${BINARY}"
echo "✓ Installed ${BINARY} to ${INSTALL_DIR}/${BINARY}"
echo ""
echo "Quick start:"
echo "  sm config set-key <your-api-key>    # Get key from https://console.supermemory.ai/keys"
echo "  sm config set-tag my-project        # Set default container tag"
echo "  sm add \"Hello, memory!\"              # Store your first memory"
echo "  sm search \"hello\"                    # Search it back"
