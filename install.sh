#!/bin/bash
set -eu

BIN_DIR=/usr/local/bin

#if [[ "$OSTYPE" != "linux-gnu"* ]] && [[ "$OSTYPE" != "darwin"* ]]; then
if [[ "$OSTYPE" != "linux-gnu"* ]]; then
  echo "dnsfsd only supports GNU/Linux"
  exit 0
fi

run_build () {
    oldwd="${PWD}"
    cd "${1}" || exit 1
    go build -o "${2}" .
    cd "${oldwd}"
}

if [ "$(id -u)" -ne 0 ]; then
    echo "You must have root permissions to run this script"
    exit 1
fi

echo "Building binaries..."

if ! run_build ./dnsfsd-util dnsfs; then
    echo "Could not build dnsfsd-util"
    exit 1
fi

if ! run_build ./daemon dnsfsd; then
    echo "Could not build daemon"
    exit 1
fi

echo "Built binaries"

if ! mv dnsfsd-util/dnsfs "$BIN_DIR" >/dev/null 2>&1; then
    echo "Could not move the built dnsfs binary to $BIN_DIR"
    exit 1
fi

if ! mv daemon/dnsfsd "$BIN_DIR" >/dev/null 2>&1; then
    echo "Could not move the built dnsfsd binary to $BIN_DIR"
    exit 1
fi

echo "Moved binaries"
echo "Install finished! Run \`dnsfs setup\` with root permission before starting dnsfsd"
