#!/bin/sh
set -eu

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

if ! mv dnsfsd-util/dnsfs /usr/local/bin/ >/dev/null 2>&1; then
    echo "Could not move the built dnsfs binary to /usr/local/bin"
    exit 1
fi

if ! mv daemon/dnsfsd /usr/local/bin/ >/dev/null 2>&1; then
    echo "Could not move the built dnsfsd binary to /usr/local/bin"
    exit 1
fi

echo "Moved binaries"
echo "Install finished! Run \`dnsfs setup\` before starting dnsfsd"
