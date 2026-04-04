#!/bin/sh
set -e

REPO="Kirskov/Digestify-My-Ci"
BINARY="digestify-my-ci"
INSTALL_DIR="/usr/local/bin"

# ── Detect OS and arch ───────────────────────────────────────────────────────

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  linux)  ;;
  darwin) ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# ── Detect package manager and install curl/tar if missing ───────────────────

install_deps() {
  if command -v curl > /dev/null 2>&1; then
    return
  fi

  echo "curl not found, installing..."

  if command -v apt-get > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq curl
  elif command -v pacman > /dev/null 2>&1; then
    pacman -Sy --noconfirm curl
  elif command -v apk > /dev/null 2>&1; then
    apk add --no-cache curl
  elif command -v dnf > /dev/null 2>&1; then
    dnf install -y curl
  elif command -v yum > /dev/null 2>&1; then
    yum install -y curl
  else
    echo "Could not install curl: no supported package manager found."
    exit 1
  fi
}

# ── Fetch latest release tag ─────────────────────────────────────────────────

latest_version() {
  curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | sed 's/.*"tag_name": *"\(.*\)".*/\1/'
}

# ── Main ─────────────────────────────────────────────────────────────────────

install_deps

VERSION="${VERSION:-$(latest_version)}"

if [ -z "$VERSION" ]; then
  echo "Could not determine latest version. Set VERSION env var to install a specific version."
  exit 1
fi

ASSET="${BINARY}-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
  ASSET="${ASSET}.exe"
fi

URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"

echo "Installing ${BINARY} ${VERSION} (${OS}/${ARCH})..."
echo "Downloading from: ${URL}"

TMP=$(mktemp)
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

# Need root to write to /usr/local/bin
if [ "$(id -u)" -eq 0 ]; then
  mv "$TMP" "${INSTALL_DIR}/${BINARY}"
elif command -v sudo > /dev/null 2>&1; then
  sudo mv "$TMP" "${INSTALL_DIR}/${BINARY}"
else
  echo "Cannot install to ${INSTALL_DIR}: not root and sudo not available."
  echo "Run as root or install manually: mv $TMP ~/bin/${BINARY}"
  exit 1
fi

echo ""
echo "${BINARY} installed to ${INSTALL_DIR}/${BINARY}"
echo "Run: ${BINARY} --help"
