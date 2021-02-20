#!/bin/bash

if [ $EUID -ne 0 ]; then
    echo "You must have root permissions to run this script"
    exit 1
fi

echo "Building binaries..."

go build -o dnsfs dnsfsd-util/main.go
if [ $? -ne 0 ]; then
    echo "Could not build dnsfsd-util"
    exit 1
fi

go build -o dnsfsd daemon
if [ $? -ne 0 ]; then
    echo "Could not build daemon"
    exit 1
fi

echo "Built binaries"

mv dnsfsd-util/dnsfs usr/local/bin/dnsfs
if [ $? -ne 0 ]; then
    echo "Could not move the built dnsfs binary to /usr/local/bin"
    exit 1
fi

mv daemon/dnsfsd /usr/local/bin/dnsfsd
if [ $? -ne 0 ]; then
    echo "Could not move the built dnsfsd binary to /usr/local/bin"
    exit 1
fi

echo "Moved binaries"
echo "Install finished! Run \`dnsfs setup\` before starting dnsfsd"