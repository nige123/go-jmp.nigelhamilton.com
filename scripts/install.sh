#!/usr/bin/env bash
set -euo pipefail

REPO="nige123/go-jmp.nigelhamilton.com"
BIN_NAME="jmp"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
REQUESTED_VERSION="${JMP_VERSION:-latest}"

fail() {
    echo "error: $*" >&2
    exit 1
}

need_cmd() {
    command -v "$1" >/dev/null 2>&1 || fail "required command not found: $1"
}

need_cmd curl
need_cmd tar
need_cmd mktemp

OS_RAW="$(uname -s)"
ARCH_RAW="$(uname -m)"

case "$OS_RAW" in
    Linux) OS="linux" ;;
    Darwin) OS="darwin" ;;
    MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
    *) fail "unsupported operating system: $OS_RAW" ;;
esac

case "$ARCH_RAW" in
    x86_64|amd64) ARCH="amd64" ;;
    i386|i486|i586|i686) ARCH="386" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) fail "unsupported architecture: $ARCH_RAW" ;;
esac

if [[ "$REQUESTED_VERSION" == "latest" ]]; then
    RELEASE_API_URL="https://api.github.com/repos/$REPO/releases/latest"
else
    RELEASE_API_URL="https://api.github.com/repos/$REPO/releases/tags/$REQUESTED_VERSION"
fi

RELEASE_JSON="$(curl -fsSL "$RELEASE_API_URL")" || fail "could not fetch release metadata"

if [[ "$OS" == "windows" ]]; then
    ASSET_PATTERN="_${OS}_${ARCH}\\.zip$"
else
    ASSET_PATTERN="_${OS}_${ARCH}\\.tar\\.gz$"
fi

ASSET_URL="$(printf "%s" "$RELEASE_JSON" | grep -Eo 'https://[^" ]+' | grep -E "$ASSET_PATTERN" | head -n1 || true)"
CHECKSUMS_URL="$(printf "%s" "$RELEASE_JSON" | grep -Eo 'https://[^" ]+' | grep -E 'checksums\\.txt$' | head -n1 || true)"

[[ -n "$ASSET_URL" ]] || fail "could not find release asset for ${OS}/${ARCH}"

TMP_DIR="$(mktemp -d)"
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

ASSET_FILE="$TMP_DIR/asset"

printf "Downloading %s\n" "$ASSET_URL"
curl -fL "$ASSET_URL" -o "$ASSET_FILE"

if [[ -n "$CHECKSUMS_URL" ]]; then
    CHECKSUMS_FILE="$TMP_DIR/checksums.txt"
    curl -fsSL "$CHECKSUMS_URL" -o "$CHECKSUMS_FILE" || fail "could not download checksums"

    ASSET_NAME="$(basename "$ASSET_URL")"
    EXPECTED_LINE="$(grep "  $ASSET_NAME" "$CHECKSUMS_FILE" || true)"
    if [[ -n "$EXPECTED_LINE" ]]; then
        EXPECTED_HASH="$(printf "%s" "$EXPECTED_LINE" | awk '{print $1}')"
        if command -v sha256sum >/dev/null 2>&1; then
            ACTUAL_HASH="$(sha256sum "$ASSET_FILE" | awk '{print $1}')"
        elif command -v shasum >/dev/null 2>&1; then
            ACTUAL_HASH="$(shasum -a 256 "$ASSET_FILE" | awk '{print $1}')"
        else
            fail "need sha256sum or shasum for checksum verification"
        fi

        [[ "$ACTUAL_HASH" == "$EXPECTED_HASH" ]] || fail "checksum mismatch for downloaded asset"
        echo "Checksum verified"
    else
        echo "warning: no checksum entry found for $ASSET_NAME"
    fi
else
    echo "warning: checksums.txt not found in release; skipping verification"
fi

EXTRACT_DIR="$TMP_DIR/extract"
mkdir -p "$EXTRACT_DIR"

if [[ "$OS" == "windows" ]]; then
    need_cmd unzip
    unzip -q "$ASSET_FILE" -d "$EXTRACT_DIR"
else
    tar -xzf "$ASSET_FILE" -C "$EXTRACT_DIR"
fi

BIN_PATH="$(find "$EXTRACT_DIR" -type f -name "$BIN_NAME" | head -n1 || true)"
if [[ -z "$BIN_PATH" && "$OS" == "windows" ]]; then
    BIN_PATH="$(find "$EXTRACT_DIR" -type f -name "${BIN_NAME}.exe" | head -n1 || true)"
fi

[[ -n "$BIN_PATH" ]] || fail "could not find extracted binary"

mkdir -p "$INSTALL_DIR"
install -m 0755 "$BIN_PATH" "$INSTALL_DIR/$BIN_NAME"

echo "Installed $BIN_NAME to $INSTALL_DIR/$BIN_NAME"
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo "Add to PATH if needed: export PATH=\"$INSTALL_DIR:\$PATH\""
fi

echo "Run: $BIN_NAME version"
