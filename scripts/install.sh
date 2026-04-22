#!/usr/bin/env sh
# CAP installer: detects OS/arch, downloads the matching binary from
# GitHub Releases and installs it into a directory on PATH.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/3110Y/cap/main/scripts/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/3110Y/cap/main/scripts/install.sh | sh -s -- --version v0.2.0
#
# Environment overrides:
#   CAP_VERSION      Tag to install (default: latest release).
#   CAP_INSTALL_DIR  Target directory (default: /usr/local/bin, or $HOME/.local/bin when not writable).
#   CAP_REPO         GitHub repo in "owner/name" form (default: 3110Y/cap).

set -eu

REPO="${CAP_REPO:-3110Y/cap}"
BIN="cap"
VERSION="${CAP_VERSION:-}"
INSTALL_DIR="${CAP_INSTALL_DIR:-}"

log()  { printf '==> %s\n' "$*"; }
warn() { printf 'warning: %s\n' "$*" >&2; }
err()  { printf 'error: %s\n' "$*" >&2; exit 1; }

while [ $# -gt 0 ]; do
    case "$1" in
        --version)
            [ $# -ge 2 ] || err "--version requires an argument"
            VERSION="$2"; shift 2 ;;
        --version=*)
            VERSION="${1#*=}"; shift ;;
        --dir)
            [ $# -ge 2 ] || err "--dir requires an argument"
            INSTALL_DIR="$2"; shift 2 ;;
        --dir=*)
            INSTALL_DIR="${1#*=}"; shift ;;
        -h|--help)
            sed -n '2,11p' "$0" | sed 's/^# \{0,1\}//'
            exit 0 ;;
        *)
            err "unknown argument: $1" ;;
    esac
done

need() {
    command -v "$1" >/dev/null 2>&1 || err "'$1' is required but not installed"
}

detect_os() {
    os="$(uname -s)"
    case "$os" in
        Linux)  printf 'linux' ;;
        Darwin) printf 'darwin' ;;
        MINGW*|MSYS*|CYGWIN*)
            err "Windows is not supported by this script; download the binary manually from https://github.com/${REPO}/releases" ;;
        *)      err "unsupported OS: $os" ;;
    esac
}

detect_arch() {
    arch="$(uname -m)"
    case "$arch" in
        x86_64|amd64)       printf 'amd64' ;;
        aarch64|arm64)      printf 'arm64' ;;
        *)                  err "unsupported architecture: $arch" ;;
    esac
}

fetch() {
    # fetch <url> <output-file>
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL -o "$2" "$1"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$2" "$1"
    else
        err "neither curl nor wget found in PATH"
    fi
}

latest_version() {
    api="https://api.github.com/repos/${REPO}/releases/latest"
    tmp="$(mktemp)"
    fetch "$api" "$tmp"
    # Extract "tag_name": "vX.Y.Z" without depending on jq.
    tag="$(sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$tmp" | head -n1)"
    rm -f "$tmp"
    [ -n "$tag" ] || err "failed to determine the latest release tag"
    printf '%s' "$tag"
}

pick_install_dir() {
    if [ -n "$INSTALL_DIR" ]; then
        printf '%s' "$INSTALL_DIR"
        return
    fi
    for candidate in /usr/local/bin "$HOME/.local/bin" "$HOME/bin"; do
        if [ -d "$candidate" ] && [ -w "$candidate" ]; then
            printf '%s' "$candidate"
            return
        fi
    done
    # Fall back to /usr/local/bin even if it requires sudo.
    printf '/usr/local/bin'
}

install_binary() {
    src="$1"
    dst_dir="$2"
    dst="$dst_dir/$BIN"

    mkdir -p "$dst_dir" 2>/dev/null || true

    if [ -w "$dst_dir" ] || { [ ! -e "$dst_dir" ] && mkdir -p "$dst_dir" 2>/dev/null; }; then
        mv "$src" "$dst"
        chmod +x "$dst"
    else
        log "Elevated privileges required to write to $dst_dir"
        need sudo
        sudo mkdir -p "$dst_dir"
        sudo mv "$src" "$dst"
        sudo chmod +x "$dst"
    fi

    printf '%s' "$dst"
}

main() {
    os="$(detect_os)"
    arch="$(detect_arch)"
    if [ -z "$VERSION" ]; then
        log "Resolving latest release of ${REPO}..."
        VERSION="$(latest_version)"
    fi
    log "Target: ${os}/${arch}, version ${VERSION}"

    asset="${BIN}_${os}_${arch}"
    url="https://github.com/${REPO}/releases/download/${VERSION}/${asset}"

    tmp_dir="$(mktemp -d)"
    trap 'rm -rf "$tmp_dir"' EXIT INT TERM
    tmp_bin="$tmp_dir/$BIN"

    log "Downloading $url"
    fetch "$url" "$tmp_bin"
    chmod +x "$tmp_bin"

    dst_dir="$(pick_install_dir)"
    log "Installing to $dst_dir/$BIN"
    dst="$(install_binary "$tmp_bin" "$dst_dir")"

    log "Installed: $dst"

    case ":${PATH}:" in
        *":${dst_dir}:"*) ;;
        *) warn "${dst_dir} is not in your PATH. Add it, for example:\n  export PATH=\"${dst_dir}:\$PATH\"" ;;
    esac

    if command -v "$BIN" >/dev/null 2>&1; then
        log "Verifying installation:"
        "$BIN" version || true
    fi
}

main "$@"
