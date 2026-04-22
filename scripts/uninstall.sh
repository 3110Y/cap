#!/usr/bin/env sh
# CAP uninstaller: removes the cap binary from the system.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/3110Y/cap/main/scripts/uninstall.sh | sh
#   sh uninstall.sh --purge       # also remove config + cache
#   sh uninstall.sh --dir ~/bin   # remove only from a specific directory
#
# Environment overrides:
#   CAP_INSTALL_DIR  Directory where the binary was installed (skips auto-detection).

set -eu

BIN="cap"
INSTALL_DIR="${CAP_INSTALL_DIR:-}"
PURGE=0

log()  { printf '==> %s\n' "$*"; }
warn() { printf 'warning: %s\n' "$*" >&2; }
err()  { printf 'error: %s\n' "$*" >&2; exit 1; }

while [ $# -gt 0 ]; do
    case "$1" in
        --dir)
            [ $# -ge 2 ] || err "--dir requires an argument"
            INSTALL_DIR="$2"; shift 2 ;;
        --dir=*)
            INSTALL_DIR="${1#*=}"; shift ;;
        --purge)
            PURGE=1; shift ;;
        -h|--help)
            sed -n '2,10p' "$0" | sed 's/^# \{0,1\}//'
            exit 0 ;;
        *)
            err "unknown argument: $1" ;;
    esac
done

remove_file() {
    target="$1"
    [ -e "$target" ] || [ -L "$target" ] || return 1

    dir="$(dirname "$target")"
    if [ -w "$dir" ]; then
        rm -f "$target"
    else
        log "Elevated privileges required to remove $target"
        command -v sudo >/dev/null 2>&1 || err "sudo not found; rerun as root or use --dir"
        sudo rm -f "$target"
    fi
    log "Removed $target"
    return 0
}

remove_dir() {
    target="$1"
    [ -e "$target" ] || return 0

    parent="$(dirname "$target")"
    if [ -w "$parent" ]; then
        rm -rf "$target"
    else
        log "Elevated privileges required to remove $target"
        command -v sudo >/dev/null 2>&1 || err "sudo not found; rerun as root"
        sudo rm -rf "$target"
    fi
    log "Removed $target"
}

main() {
    removed=0

    if [ -n "$INSTALL_DIR" ]; then
        candidates="$INSTALL_DIR/$BIN"
    else
        candidates="/usr/local/bin/$BIN $HOME/.local/bin/$BIN $HOME/bin/$BIN /usr/bin/$BIN /opt/homebrew/bin/$BIN"
        # Also honor whatever `cap` resolves to on the current PATH.
        if command -v "$BIN" >/dev/null 2>&1; then
            resolved="$(command -v "$BIN")"
            case " $candidates " in
                *" $resolved "*) ;;
                *) candidates="$candidates $resolved" ;;
            esac
        fi
    fi

    for path in $candidates; do
        if [ -e "$path" ] || [ -L "$path" ]; then
            remove_file "$path" && removed=$((removed + 1)) || true
        fi
    done

    if [ "$removed" -eq 0 ]; then
        warn "cap binary not found in the usual locations; nothing to remove"
    fi

    if [ "$PURGE" -eq 1 ]; then
        log "Purging CAP config and cache"
        remove_dir "$HOME/.config/cap"
        remove_dir "$HOME/.cache/cap"
    else
        if [ -d "$HOME/.config/cap" ] || [ -d "$HOME/.cache/cap" ]; then
            log "Config/cache kept. Use --purge to also remove:"
            [ -d "$HOME/.config/cap" ] && printf '      %s\n' "$HOME/.config/cap"
            [ -d "$HOME/.cache/cap" ]  && printf '      %s\n' "$HOME/.cache/cap"
        fi
    fi

    if command -v "$BIN" >/dev/null 2>&1; then
        warn "'$BIN' is still reachable at $(command -v "$BIN"). You may have another copy installed."
    else
        log "Uninstall complete"
    fi
}

main "$@"
