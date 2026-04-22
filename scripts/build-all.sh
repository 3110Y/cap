#!/usr/bin/env bash
set -euo pipefail

BIN="cap"
OUT_DIR="bin"

mkdir -p "${OUT_DIR}"

platforms=(
    "darwin amd64"
    "darwin arm64"
    "linux amd64"
    "linux arm64"
    "windows amd64"
    "windows arm64"
)

for p in "${platforms[@]}"; do
    os="${p% *}"
    arch="${p#* }"
    ext=""
    if [ "${os}" = "windows" ]; then
        ext=".exe"
    fi
    output="${OUT_DIR}/${BIN}_${os}_${arch}${ext}"
    echo "Building ${output}..."
    GOOS="${os}" GOARCH="${arch}" CGO_ENABLED=0 go build \
        -ldflags "-s -w" \
        -o "${output}" \
        ./cmd/"${BIN}"
done

echo "Done. Artifacts in ${OUT_DIR}/"
